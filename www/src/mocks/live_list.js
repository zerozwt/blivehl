const LiveList = {
    code: 0,
    msg: "",
    data: {
        list: [
            {
                title: "直播标题",
                live_start_time: "2023-04-02 20:00:00 CST",
                cover: "/icon.png",
                live_id: 123456789,
            },
            {
                title: "上一场直播",
                live_start_time: "2023-04-01 20:00:00 CST",
                cover: "/icon.png",
                live_id: 113456789,
            },
            {
                title: "再上一场直播",
                live_start_time: "2023-03-30 20:00:00 CST",
                cover: "/icon.png",
                live_id: 103456789,
            },
        ],
        ended: false
    }
}

const AdminLiveList = (opt) => {
    let ret = {
        code: 0,
        msg: "",
        data: {list: [], total: 100},
    }
    for (let i = 0; i < 10; i++) {
        ret.data.list.push({
            title: `直播${i+1}`,
            live_start_time: "2023-03-30 20:00:00 CST",
            cover: "/icon.png",
            live_id: 100-i
        })
    }
    return ret
}

export default {
    'get|^/api/live/list': opt => LiveList,
    'get|^/api/admin/lives': AdminLiveList
}