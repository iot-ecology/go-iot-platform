import { scriptWaringById } from '@/services/ant-design-pro/api';
import React, { useEffect, useState } from 'react';

export type Props = {
  id: string;
};

const ScriptNameShow: React.FC<Props> = (props) => {
  const [data, setData] = useState<string>('');
  useEffect(() => {
    const fetchData = async () => {
      const response = await scriptWaringById(props.id);
      setData(response.data.name);
    };

    fetchData();

    // 清理函数，用于取消异步操作或清理副作用
    return () => {
      // 如果有需要取消请求的逻辑，可以在这里实现
    };
  }, [props.id]);

  return (
    <>
      <div>{data || '-'}</div>
    </>
  );
};

export default ScriptNameShow;
