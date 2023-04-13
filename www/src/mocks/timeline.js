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

export default {
    'get|^/api/highlight/timeline': opt => Timeline
}