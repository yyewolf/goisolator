# Put random value in "test"
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 5 | head -n 1 > test

docker build -t localhost:5000/test:latest .
docker push localhost:5000/test:latest