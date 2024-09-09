import React from "react";



export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.UserListItem) => void;
  onSubmit: (values: API.UserListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.UserListItem;
};
const UserUpdateForm: React.FC<UpdateFormProps> = (props) => {
  return (
   <div>
     a1
   </div>
  );
};

 export default  UserUpdateForm;
