# $prefixarr defines which strategies and periodicities to run experiments on.
# so "mfe_5" would be strategy=most frequently encountered, with periodicity=5
$prefixarr = `
    "rs_5", "rs_10", "rs_15"
    # "mfe_5", "mfe_10", "mfe_15", `
    # "mru_5", "mru_10", "mru_15", `
    

foreach ($prefix in $prefixarr) {
    # Start main.go server
    Write-Output "$($prefix): Starting main.go server"
    $proc = Start-Process .\main.exe -ArgumentList ".\confs\$prefix.yaml", '.\lru_config.yaml', '.\functions_test' -PassThru -WindowStyle normal
    #   -RedirectStandardOutput ".\exp_out\z_$prefix.go_server_console.out" -RedirectStandardError ".\exp_out\z_$prefix.go_server_console.err"
    Write-Output $proc.Id

    Start-Sleep -Seconds 10
    # Start python simulate.py
    Write-Output "$($prefix): Starting python simulate.py"
    Start-Process python -ArgumentList 'simulate.py', 'test_workload', 'functions_test', 'http://127.0.0.1:1024' -WindowStyle normal -Wait
        # -Wait -RedirectStandardOutput ".\exp_out\z_$prefix.py_sim_console.out" -RedirectStandardError ".\exp_out\z_$prefix.py_sim_console.err"

    # This dumps info to taskmaster_activation_ids.txt
    Start-Sleep -Seconds 3
    Write-Output "$($prefix): Dumping data..."
    Invoke-RestMethod -Uri "127.0.0.1:1024/dumpData" -Method GET
    Start-Sleep -Seconds 7

    # Copy and rename the openwhisk log file
    Write-Output "$($prefix): Copying openwhisk_out"
    Copy-Item "..\..\openwhisk\openwhisk_out" -Destination ".\exp_out\$($prefix)_openwhisk_out"

    # Stop main.go server
    Write-Output "$($prefix): Stopping main.go server"
    Stop-Process -Id $proc.Id
    Start-Sleep -Seconds 20
}
