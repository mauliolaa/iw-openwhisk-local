$prefixarr = "rs_5", "rs_10", "rs_15", "mfe_5", "mfe_10", "mfe_15"
foreach ($prefix in $prefixarr) {

$func_stats = @{
    "hello_js"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "fact_js"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "simple_math"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "fact_py"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "hello_ruby"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "fact_ruby"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "hello_php"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "fact_php"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "hello_java"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    "fact_java"=@{"prewarmedContainerCount"=0;        
    "coldContainerCount"=0;
    "warmedContainerCount"=0;
    "recreatedContainerCount"=0};
    
    }
foreach($line in [System.IO.File]::ReadLines("taskmaster\exp_out\$($prefix)_openwhisk_out")) {
    foreach($key in $func_stats.Keys) {
        if ($line -match $key) {
            if ($line -match "containerState: prewarmed container") {
                $func_stats[$key]["prewarmedContainerCount"] += 1
            }
            elseif ($line -match "containerState: cold container") {
                $func_stats[$key]["coldContainerCount"] += 1
            }
            elseif ($line -match "containerState: warmed container") {
                $func_stats[$key]["warmedContainerCount"] += 1
            }
            elseif ($line -match "containerState: recreated container") {
                $func_stats[$key]["recreatedContainerCount"] += 1
            }
        }
    }
}
# foreach value in $func_stats.Values()
# Write-Host $func_stats.GetEnumerator() | sort value -descending
$json = Get-Content -Path "$($prefix)_results.json" -Raw | ConvertFrom-Json

foreach ($h in $func_stats.GetEnumerator() )
{
    # Write-Host "$($h.Name) :"
    foreach ($j in $h.Value.GetEnumerator() ) {
        # Write-Host "$($j.Name) : $($j.Value)"
        $json.$($h.Name).$($j.Name) = $j.Value
    }
}

$json = $json | ConvertTo-Json -Compress

# Save JSON to file
$json | Set-Content -Path "$($prefix)_results_fixed.json"

}