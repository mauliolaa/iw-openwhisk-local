$prefixarr = "rs_10", "rs_15", "mfe_5", "mfe_10", "mfe_15"
# "rs_5", "rs_10", "rs_15"
# python .\get_experiment_metrics.py .\taskmaster\exp_out\diff\mfe_5_openwhisk_out .\taskmaster\functions_test http://127.0.0.1:1024 .\taskmaster\mfe_5_taskmaster_activation_ids.txt
$previous = ".\_openwhisk_out" 
foreach ($prefix in $prefixarr) {
    $current = ".\$($prefix)_openwhisk_out"
    $dest = ".\diff\$($prefix)_openwhisk_out"
    (Compare-Object -ReferenceObject (Get-Content $previous) -DifferenceObject (Get-Content $current)).InputObject > $dest
    $previous = $current
}
