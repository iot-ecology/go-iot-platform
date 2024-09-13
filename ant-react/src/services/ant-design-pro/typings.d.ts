// @ts-ignore
/* eslint-disable */

declare namespace API {
  type CurrentUser = {
    name?: string;
    avatar?: string;
    userid?: string;
    email?: string;
    signature?: string;
    title?: string;
    group?: string;
    tags?: { key?: string; label?: string }[];
    notifyCount?: number;
    unreadCount?: number;
    country?: string;
    access?: string;
    geographic?: {
      province?: { label?: string; key?: string };
      city?: { label?: string; key?: string };
    };
    address?: string;
    phone?: string;
  };

  type LoginResult = {
    status?: string;
    type?: string;
    currentAuthority?: string;
  };

  type PageParams = {
    current?: number;
    pageSize?: number;
  };

  type RuleListItem = {
    key?: number;
    disabled?: boolean;
    href?: string;
    avatar?: string;
    name?: string;
    owner?: string;
    desc?: string;
    callNo?: number;
    status?: number;
    updatedAt?: string;
    createdAt?: string;
    progress?: number;
  };

  type UserListItem = {
    ID?: number;
    username?: string;
    password?: string;
    email?: string;
  };
  type SimListItem = {
    ID?: number;
    access_number?: string;
    iccid?: string;
    imsi?: string;
    operator?: string;
    expiration?: string;
  };
  type FeiShuListItem = {
    ID?: number;
    name?: string;
    access_token?: string;
    secret?: string;
    content?: string;
  };
  type DingDingListItem = {
    ID?: number;
    name?: string;
    access_token?: string;
    secret?: string;
    content?: string;
  };
  type SignalListItem = {
    ID?: number;
    protocol?: string;
    identification_code?: string;
    device_uid?: string;
    name?: string;
    alias?: string;
    type?: string;
    unit?: string;
    cache_size?: string;
  };
  type MqttListItem = {
    ID?: number;
    host?: string;
    port?: number;
    client_id?: string;
    username?: string;
    password?: string;
    subtopic?: string;
    start?: boolean;
    last_push_time?: string;
    script?: string;
  };
  type RecordInfo = {
    product_id: ?number;
    quantity?: number;
  };

  type ShipmentRecordListItem = {
    ID?: number;
    shipment_date?: string;
    technician?: string;
    customer_name?: string;
    customer_phone?: string;
    customer_address?: string;
    tracking_number?: string;
    status?: string;
    description?: string;

    product_plans?: ProductPlanCreateParam[] | [];
  };

  type ProductPlanCreateParam = {
    product_id?: number;
    quantity?: number;
    shipment_record_id?: number;
    device_info_id?: number;
  };
  type DeviceGroupItem = {
    ID?: number;
    name?: string;
  };
  type ProductItem = {
    ID?: number;
    name?: string;
    description?: string;
    sku?: string;
    price?: number;
    cost?: number;
    quantity?: number;
    minimum_stock?: number;
    warranty_period?: number;
    status?: string;
    tags?: string;
    image_url?: string | FileItemUpload[];
  };
  type FileItemUpload = {
    response: FileUpdateResponse;
  };
  type FileUpdateResponse = {
    file_path: string;
    message: string;
  };
  type DeptListItem = {
    ID?: number;
    name?: string;
    parent_id?: number;
    parent_name?: string;
  };
  type DeviceInfoItem = {
    ID?: number;
    product_id?: number;
    sn?: string;
    manufacturing_date?: string;
    procurement_date?: string;
    source?: number;
    warranty_expiry?: string;
    push_interval?: string;
    error_rate?: number;
    protocol?: string;
  };
  type RoleListItem = {
    ID?: number;
    name?: string;
    description?: string;
  };
  type MessageListItem = {
    ID?: number;
    content?: string;
    en_content?: string;
    message_type_id?: number;
  };

  type RuleList = {
    data?: RuleListItem[] /** 列表的内容总数 */;
    total?: number;
    success?: boolean;
  };

  type PageList<T> = {
    data?: T[];
    total?: number;
    page?: number;
    size?: number;
  };
  type UserList = {
    data?: PageList<UserListItem> /** 列表的内容总数 */;
    message?: string;
    code?: number;
  };
  type CommonPage<T> = {
    data?: PageList<T> /** 列表的内容总数 */;
    message?: string;
    code?: number;
  };

  type CommonResp<T> = {
    data?: T;
    message?: string;
    code?: number;
  };

  type FakeCaptcha = {
    code?: number;
    status?: string;
  };

  type LoginParams = {
    username?: string;
    password?: string;
    autoLogin?: boolean;
    type?: string;
  };

  type ErrorResponse = {
    /** 业务约定的错误码 */
    errorCode: string /** 业务上的错误信息 */;
    errorMessage?: string /** 业务上的请求是否成功 */;
    success?: boolean;
  };

  type NoticeIconList = {
    data?: NoticeIconItem[] /** 列表的内容总数 */;
    total?: number;
    success?: boolean;
  };

  type NoticeIconItemType = 'notification' | 'message' | 'event';

  type NoticeIconItem = {
    id?: string;
    extra?: string;
    key?: string;
    read?: boolean;
    avatar?: string;
    title?: string;
    status?: string;
    datetime?: string;
    description?: string;
    type?: NoticeIconItemType;
  };
}
