dbAdmin = db.getSiblingDB("admin");
dbAdmin.createUser({
    user: "iot",
    pwd: "iot123",
    roles: [{ role: "userAdminAnyDatabase", db: "admin" }],
    mechanisms: ["SCRAM-SHA-1"],
});

// Authenticate user
dbAdmin.auth({
    user: "iot",
    pwd: "iot123",
    mechanisms: ["SCRAM-SHA-1"],
    digestPassword: true,
});

use iot;
db.createCollection("calc");
db.createCollection("waring");
db.createCollection("script_waring");

// todo 导入数据可以在此处操作,或者新建一个单独的用户来操作有权限的库
