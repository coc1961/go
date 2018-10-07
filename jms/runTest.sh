docker run -d --name='activemq' -it --rm -e 'ACTIVEMQ_CONFIG_MINMEMORY=512' -e 'ACTIVEMQ_CONFIG_MAXMEMORY=2048' -p 1883:1883/tcp -p 5672:5672/tcp -p 8161:8161/tcp -p 61613:61613/tcp -p 61614:61614/tcp -p 61616:61616/tcp webcenter/activemq:latest
./runAllTest.sh
docker stop activemq


