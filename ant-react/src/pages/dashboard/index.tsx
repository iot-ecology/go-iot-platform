import DeviceStatsChart from "@/pages/dashboard/DeviceStatsChart";
import ProtocolStatsCards from "@/pages/dashboard/ProtocolStatsTable";
import MessageList from "@/pages/dashboard/MessageList";
import {Col, Row} from "antd";
import DeviceMap from "@/pages/dashboard/DeviceMap";
import ProductDeviceChart from "@/pages/dashboard/ProductDeviceCards";
import AlarmStatistics from "@/pages/dashboard/AlarmStatistics";

const Dashboard = () => {
  return (
    <>
      <DeviceMap />
      <Row gutter={16}>
        <Col span={12}>
          <DeviceStatsChart />
        </Col>
        <Col span={12}>
          <ProtocolStatsCards />
        </Col>
      </Row>
      <Row gutter={16}>
        <Col span={8}>
          <AlarmStatistics />
        </Col>
        <Col span={8}>
          <MessageList />
        </Col>
        <Col span={8}>
          <ProductDeviceChart />
        </Col>
      </Row>
    </>
  );
};
export default Dashboard
