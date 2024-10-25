import React, { useEffect, useState } from 'react';
import { Card, Col, Row, Modal } from 'antd';
import axios from 'axios';

const MqttMirror = () => {
  const [mqttData, setMqttData] = useState([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [modalData, setModalData] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:8080/pod_mqtt');
        setMqttData(response.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, []);

  const handleCardClick = async (nodeName) => {
    try {
      const response = await axios.get(`http://localhost:8080/mqtt/node-using-status?name=${nodeName}`);
      const data = JSON.parse(response.data.data); // 解析数据
      setModalData(data);
      setModalVisible(true);
    } catch (error) {
      console.error('Error fetching node status:', error);
    }
  };

  const handleCloseModal = () => {
    setModalVisible(false);
    setModalData(null);
  };

  return (
    <div>
      <h2>MQTT Nodes</h2>
      <Row gutter={16}>
        {mqttData.map((pod) => (
          <Col span={8} key={pod.name}>
            <Card title={pod.name} bordered={false} onClick={() => handleCardClick(pod.name)}>
              <p>Host: {pod.host}</p>
              <p>Port: {pod.port}</p>
              <p>Type: {pod.type}</p>
              <p>Size: {pod.size} 个</p>
            </Card>
          </Col>
        ))}
      </Row>
      <Modal
        title="Node Status"
        visible={modalVisible}
        onCancel={handleCloseModal}
        footer={null}
      >
        {modalData ? (
          <div>
            {modalData.data.map((client) => (
              <Card key={client.client_ids[0]}
                    style={{marginBottom: 16, border: '1px solid #e0e0e0', borderRadius: '8px', padding: '16px'}}>


                {client.client_infos.map((info, index) => (
                  <div key={index}
                       style={{padding: '8px', border: '1px solid #f0f0f0', borderRadius: '4px', marginBottom: 8}}>
                    <h5 style={{margin: '0 0 4px'}}>Client Info {index + 1}</h5>
                    <p style={{margin: '4px 0'}}><strong>Broker:</strong> {info.broker}</p>
                    <p style={{margin: '4px 0'}}><strong>Port:</strong> {info.port}</p>
                    <p style={{margin: '4px 0'}}><strong>Username:</strong> {info.username}</p>
                    <p style={{margin: '4px 0'}}><strong>Password:</strong> {info.password}</p>
                    <p style={{margin: '4px 0'}}><strong>Sub Topic:</strong> {info.sub_topic}</p>
                    <p style={{margin: '4px 0'}}><strong>Client ID:</strong> {info.client_id}</p>
                  </div>
                ))}
              </Card>
            ))}
          </div>
        ) : (
          <p>Loading...</p>
        )}
      </Modal>
    </div>
  );
};

export default MqttMirror;
