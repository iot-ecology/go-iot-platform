function createMqttClient(i) {
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    const raw = JSON.stringify({
        "client_id": "TT_" + i ,
        "host": "127.0.0.1",
        "port": 1883,
        "username": "admin",
        "password": "admin",
        "subtopic": "/test_topic/" + i
    });

    const requestOptions = {
        method: "POST",
        headers: myHeaders,
        body: raw,
        redirect: "follow"
    };

    fetch("http://127.0.0.1:8005/mqtt/create", requestOptions)
        .then((response) => response.text())
        .then((result) => console.log(result))
        .catch((error) => console.error(error));
}

function callCreateMqtt() {
    for (let i = 0; i < 100; i++) {

        createMqttClient(i);

    }


}


function setScript(i ){
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    const raw = JSON.stringify({
        "id": i,
        "script": "function main(nc) {\n   var result = JSON.parse(nc);\n    return [result];\n}"
    });

    const requestOptions = {
        method: "POST",
        headers: myHeaders,
        body: raw,
        redirect: "follow"
    };

    fetch("http://127.0.0.1:8005/mqtt/set-script", requestOptions)
        .then((response) => response.text())
        .then((result) => console.log(result))
        .catch((error) => console.error(error));
}


function callSetScript(){
    for (let i = 0; i < 100; i++) {

        setScript(i );

    }
}







function c(mqtt_client_id, sid) {
    var now = new Date();
    var formattedNow = now.toISOString().slice(0, 19).replace('T', ' ');

    var cc = `INSERT INTO \`signals\` (\`protocol\`, \`identification_code\`, \`device_uid\`, \`name\`, \`alias\`, \`type\`, \`unit\`, \`cache_size\`, \`created_at\`, \`updated_at\`, \`deleted_at\`) 
              VALUES ('mqtt', '${mqtt_client_id}', ${mqtt_client_id}, '信号-${sid}', '信号-${sid}', '数字', '无', 1, '${formattedNow}', '${formattedNow}', NULL);`;
    console.log(cc);
}

function main(){
    for (let i = 0; i <100; i++) {
        for (let j = 0; j <200; j++) {
            c(i+1 , j);
        }
    }
}
// callCreateMqtt()
// callSetScript()
main()