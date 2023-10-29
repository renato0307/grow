BINARY_NAME=monitorghm
DESTINATION=growzero1
SERVICE=${BINARY_NAME}.service

build:
	env GOOS=linux GOARCH=arm GOARM=6 go build -o ./bin/${BINARY_NAME}
	
publish:
	scp ./bin/${BINARY_NAME} ${DESTINATION}:~/
	ssh ${DESTINATION} 'sudo mv ${BINARY_NAME} /usr/local/bin'

install-service:
	scp ./${SERVICE} ${DESTINATION}:~/
	ssh ${DESTINATION} 'sudo id -u ${BINARY_NAME} &>/dev/null || useradd -s /sbin/nologin -M ${BINARY_NAME}'
	ssh ${DESTINATION} 'sudo chmod a+rw /dev/gpiochip0'
	ssh ${DESTINATION} 'sudo mv ${SERVICE} /etc/systemd/system/'
	ssh ${DESTINATION} 'sudo chmod 755 /etc/systemd/system/${SERVICE}'
	ssh ${DESTINATION} 'sudo systemctl daemon-reload'
	ssh ${DESTINATION} 'sudo systemctl enable ${SERVICE}'
	ssh ${DESTINATION} 'sudo systemctl restart ${SERVICE}'

check-service-logs:
	ssh ${DESTINATION} 'sudo journalctl -u ${SERVICE}'