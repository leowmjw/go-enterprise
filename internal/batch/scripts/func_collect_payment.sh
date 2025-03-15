function myip() {
	curl http://icanhazip.com

	ip addr | grep inet$IP | \
	cut -d"/" -f 1 | \
	grep -v 127\.0 | \
	grep -v \:\:1 | \
	awk '{$1=$1};1'
}