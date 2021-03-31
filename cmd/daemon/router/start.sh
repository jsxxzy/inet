# go build -o inetd ..
scp inetd admin@10.32.0.1:/tmp

cat ./run.sh | ssh admin@10.32.0.1
