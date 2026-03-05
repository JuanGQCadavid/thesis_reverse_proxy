KIND ?= DEFAULT
NOW := $(shell date '+%F_%H:%M:%S')
WORK_DIR := ${KIND}/${NOW}

capture:
	mkdir -p ${WORK_DIR}
	echo ${WORK_DIR}
	touch ${WORK_DIR}/c1.pcap
	tshark -i lo \
	  -f "tcp port 9000 or udp port 9000" \
	  -w "${WORK_DIR}/c1.pcap"

PATH?=DEFAULT
convert:
	mkdir -p ${PATH}/output
	tcptrace --output_dir="${PATH}/output" -l -e ${PATH}/c1.pcap 
