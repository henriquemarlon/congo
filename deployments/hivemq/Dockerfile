FROM hivemq/hivemq4:latest

COPY --chown=hivemq:hivemq ./deployments/hivemq/kafka-extension/config.xml /opt/hivemq/extensions/hivemq-kafka-extension/conf/config.xml

RUN rm -f /opt/hivemq/extensions/hivemq-kafka-extension/DISABLED