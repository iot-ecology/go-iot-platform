import React from 'react';
import { Tabs } from 'antd';
import { CoapMirror } from "@/pages/mirror/protocol/coap";
import { HttpMirror } from "@/pages/mirror/protocol/http";
import  MqttMirror from "@/pages/mirror/protocol/mqtt";
import { TcpMirror } from "@/pages/mirror/protocol/tcp";
import { WsMirror } from "@/pages/mirror/protocol/ws";

const { TabPane } = Tabs;

const F = () => {

  return (

    <Tabs
      defaultActiveKey="1"
      tabBarGutter={20}
    >
      <TabPane tab="CoAP" key="1">
        <CoapMirror />
      </TabPane>
      <TabPane tab="HTTP" key="2">
        <HttpMirror />
      </TabPane>
      <TabPane tab="MQTT" key="3">
        <MqttMirror />
      </TabPane>
      <TabPane tab="TCP" key="4">
        <TcpMirror />
      </TabPane>
      <TabPane tab="WebSocket" key="5">
        <WsMirror />
      </TabPane>
    </Tabs>
  );
};

export default F;
