const Prepare = {
    code: 0,
    msg: "",
    data: {
        live_status: 1,
        title: "直播标题",
        live_start_time: "2023-04-02 20:00:00 CST",
        cover: "/icon.png",
        live_id: 123456789,
        light_ts: 1234,
    }
}

export default {
    'get|^/api/prepare': opt => Prepare
}