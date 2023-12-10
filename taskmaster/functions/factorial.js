// name: str
// place: str
function factorial(n) {
    if (n <= 1) {
        return n;
    }
    return n * factorial(n-1);
}

function main(params) {
    return {payload: factorial(parseInt(params.n))}
}