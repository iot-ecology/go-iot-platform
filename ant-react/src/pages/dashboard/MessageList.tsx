import React, { useState } from 'react';
import { Card, Tabs, List, Typography } from 'antd';

const { Title } = Typography;
const { TabPane } = Tabs;

const messages = {
  "分类1": [
    "消息内容 1-1: 这是第一类的第一条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    "消息内容 1-2: 这是第一类的第二条消息。",
    // ...其他消息
    "消息内容 1-3: 这是第一类的第三条消息。",
  ],
  "分类2": [
    "消息内容 2-1: 这是第二类的第一条消息。",
    "消息内容 2-2: 这是第二类的第二条消息。",
    "消息内容 2-3: 这是第二类的第三条消息。",
  ],
  "分类3": [
    "消息内容 3-1: 这是第三类的第一条消息。",
    "消息内容 3-2: 这是第三类的第二条消息。",
    // ...其他消息
  ],
};

const MessageCard: React.FC = () => {
  const [activeKey, setActiveKey] = useState("分类1");

  return (
    <Card >
      <Title level={4}>消息列表</Title>
      <Tabs activeKey={activeKey} onChange={setActiveKey}>
        {Object.keys(messages).map(category => (
          <TabPane tab={category} key={category}>
            <List
              bordered
              dataSource={messages[category]}
              renderItem={item => (
                <List.Item>
                  {item}
                </List.Item>
              )}
              style={{ maxHeight: 200, overflowY: 'auto' }} // 设置最大高度和滚动条
            />
          </TabPane>
        ))}
      </Tabs>
    </Card>
  );
};

export default MessageCard;
