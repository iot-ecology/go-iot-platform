FROM rabbitmq:3.10.6-management

ADD plugins/rabbitmq_delayed_message_exchange-3.10.2.ez /plugins

# 开启插件
RUN rabbitmq-plugins enable rabbitmq_delayed_message_exchange

ENTRYPOINT ["rabbitmq-server"]