<?php

function factorial($n) {
    if ($n === 0 || $n === 1) {
        return 1;
    } else {
        return $n * factorial($n - 1);
    }
}

function main(array $args) {
    $result = factorial($args["n"]);
    return ["result" => $result];
}