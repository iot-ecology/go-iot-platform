import { mqttById } from '@/services/ant-design-pro/api';
import React, { useEffect, useState } from 'react';

export type Props = {
  protocol: string;
  device_uid: string;
};

const DeviceUidShow: React.FC<Props> = (props) => {
  const [data, setData] = useState<string>('');
  useEffect(() => {
    const fetchData = async () => {
      try {
        if (props.protocol === 'MQTT') {
          const response = await mqttById(props.device_uid);
          setData(response.data.client_id);
        } else {
          setData('-');
        }
      } catch (error) {
        console.error('Error fetching MQTT client ID:', error);
        setData('-');
      }
    };

    fetchData();

    // 清理函数，用于取消异步操作或清理副作用
    return () => {
      // 如果有需要取消请求的逻辑，可以在这里实现
    };
  }, [props.protocol, props.device_uid]);

  return (
    <>
      <div>{data || '-'}</div>
    </>
  );
};

export default DeviceUidShow;
