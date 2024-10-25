import React, { useEffect, useState } from 'react';
import { Card, Col, Row } from 'antd';
import { podInfo } from "@/services/ant-design-pro/api";
import PodInfoModal from "@/pages/mirror/protocol/MetricsModal";

const WsMirror = () => {
  const [coapData, setCoapData] = useState([]);
  const [selectedPod, setSelectedPod] = useState(null);
  const [modalVisible, setModalVisible] = useState(false);

  useEffect(() => {
    // 请求数据
    const fetchData = async () => {
      try {
        const response = await podInfo();
        const data = response;

        // 提取 CoAP 数据
        const coapInfo = data.find(item => item.type === 'ws');
        if (coapInfo) {
          setCoapData(coapInfo.pod_info);
        }
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, []);

  const handleCardClick = (pod) => {
    setSelectedPod(pod);
    setModalVisible(true);
  };

  const handleCloseModal = () => {
    setModalVisible(false);
    setSelectedPod(null);
  };

  return (
    <div>
      <h2>WebSocket Nodes</h2>
      <Row gutter={16}>
        {coapData.map((pod, index) => (
          <Col span={8} key={index}>
            <Card title={pod.name} bordered={false} onClick={() => handleCardClick(pod)}>
              <p>Host: {pod.node_info.host}</p>
              <p>Port: {pod.node_info.port}</p>
            </Card>
          </Col>
        ))}
      </Row>
      {selectedPod && (
        <PodInfoModal

          visible={modalVisible}
          pod={selectedPod}
          type={"ws"}
          onClose={handleCloseModal}
        />
      )}
    </div>
  );
};

export { WsMirror };
