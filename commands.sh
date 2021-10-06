set -ex
run(){
  go install
  req_until_slow -url $1 -timeout 1s
}
$@