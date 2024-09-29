import DeviceMap from "@/pages/dashboard/DeviceMap";
import ProtocolStatsCards from "@/pages/dashboard/ProtocolStatsTable";
import MessageList from "@/pages/dashboard/MessageList";
import { Col, Row } from "antd";
import AlarmStatistics from "@/pages/dashboard/AlarmStatistics";
import ProductDeviceChart from "@/pages/dashboard/ProductDeviceCards";

const Dashboard = () => {
  return (
    <>
      <Row gutter={[16,16]} style={{height:"auto"}}>
        <Col span={12} style={{height:"700px"}}>
          <DeviceMap />
        </Col>
        <Col span={12}  >
          <ProtocolStatsCards />
          <AlarmStatistics  />
          <MessageList />
        </Col>
      </Row>

    </>
  );
};

export default Dashboard;
