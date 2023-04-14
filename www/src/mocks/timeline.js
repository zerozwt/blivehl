const Timeline = {
    code: 0,
    msg: "",
    data: {
        timeline: [
            {time: 5542, comment: "这一段绝了"},
            {time: 4321, comment: "哈哈哈"},
            {time: 3456, comment: "这里可以切"},
            {time: 2345, comment: "这个活整的好"},
            {time: 1234, comment: "此处高能"},
        ]
    }
}

const AdminTimeline = (opt) => {
    let ret = {
        code: 0,
        msg: "",
        data: {
            timeline: [],
        },
    }
    for (let i = 0; i < 20; i++) {
        ret.data.timeline.push({
            time: 3000+60*i,
            comment: "哈哈哈",
            author: `点灯员${(i%4)+1}`
        })
    }
    return ret
}

export default {
    'get|^/api/highlight/timeline': opt => Timeline,
    'get|^/api/admin/timeline': AdminTimeline,
}