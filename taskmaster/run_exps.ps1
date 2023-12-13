# prefix of current iteration
# $prefixarr = "lru_p1", "lru_p5", "lru_p10", `
#     "mfe_p1", "mfe_p5", "mfe_p10", `
#     "mru_p1", "mru_p5", "mru_p10", `
#     "pq_p1", "pq_p5", "pq_p10", `
#     "rs_p1", "rs_p5", "rs_p10"
$prefixarr = 
    "mfe_5", "mfe_10", "mfe_15", `
    "mru_5", "mru_10", "mru_15", `
    "rs_5", "rs_10", "rs_15"

foreach ($prefix in $prefixarr) {
    # Start main.go server
    Write-Output "Starting main.go server"
    $proc = Start-Process go -ArgumentList 'run', 'main.go', ".\confs\$prefix.yaml", '.\lru_config.yaml', '.\functions_test' `
      -RedirectStandardOutput ".\exp_out\z_$prefix.go_server_console.out" -RedirectStandardError ".\exp_out\z_$prefix.go_server_console.err" -PassThru

    Start-Sleep -Seconds 10
    # Start python simulate.py
    Write-Output "Starting python simulate.py"
    Start-Process python -ArgumentList 'simulate.py', 'test_workload_short', 'functions_test', 'http://127.0.0.1:1024' `
        -Wait -RedirectStandardOutput ".\exp_out\z_$prefix.py_sim_console.out" -RedirectStandardError ".\exp_out\z_$prefix.py_sim_console.err"

    # This dumps info to taskmaster_activation_ids.txt
    Invoke-RestMethod -Uri "127.0.0.1:1024/dumpData" -Method GET
    Start-Sleep -Seconds 10

    Start-Sleep -Seconds 7
    # Copy and rename the openwhisk log file
    Copy-Item "..\..\openwhisk\openwhisk_out" -Destination ".\exp_out\$($prefix)_openwhisk_out"

    # Stop main.go server
    Write-Output "Stopping main.go server"
    Stop-Process -InputObject $proc
    Start-Sleep -Seconds 10
}
