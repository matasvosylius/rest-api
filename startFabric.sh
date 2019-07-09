#!/bin/bash


function dkcl(){
        CONTAINER_ID=$(docker ps -aq)
	echo
        if [ -z "$CONTAINER_ID" -o "$CONTAINER_ID" = " " ]; then
                echo "========== You're already fresh and clean... There are no containers available for deletion =========="
        else
                docker rm -f $CONTAINER_ID
        fi
	echo
}

function dkrm(){
        DOCKER_IMAGE_ID=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
	echo
        if [ -z "$DOCKER_IMAGE_ID" -o "$DOCKER_IMAGE_ID" = " " ]; then
		echo "========= You're already flesh and clean... No images available for deletion ==========="
        else
                docker rmi -f $DOCKER_IMAGE_ID
        fi
	echo
}

function restartNetwork() {
	echo

  #teardown the network and clean the containers and intermediate images
	docker-compose -f ./artifacts/docker-compose.yaml down
	dkcl
	dkrm

	#Cleanup the stores
	rm -rf ./fabric-client-kv-org*

	#Start the network
	docker-compose -f ./artifacts/docker-compose.yaml up -d
	echo
}


restartNetwork
