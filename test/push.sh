# Put random value in "test"
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 5 | head -n 1 > test

docker build -t test .