function setPage(id, pages, page_no, op) {
    // id: 页签父元素; pages: 后端总页数; page_no: 获取当前页码; op: 操作[+-]
    // 页左侧按钮
    if (page_no != 1) {
        $(id + ">.btn-prev").attr("disabled", false)
    } else {
        $(id + ">.btn-prev").attr("disabled", true)
    }
    // 页右侧按钮
    if (page_no != pages) {
        $(id + ">.btn-next").attr("disabled", false)
    } else {
        $(id + ">.btn-next").attr("disabled", true)
    }
    // 页面渲染
    if (pages == 1) {
        var p = "<li class='number active'>1</li>"
        $(id + ">.el-pager").html(p)
    } else if (pages <= 5) {
        // 总页数 <= 5
        var p = ""
        for (var i = 1; i <= pages; i++) {
            if (page_no == i) {
                p += "<li class='number active'>" + i + "</li>"
            } else {
                p += "<li class='number'>" + i + "</li>"
            }
        }
        $(id + ">.el-pager").html(p)
    } else {
        // 总页数 > 5
        // 左移或右移到两端时需要重新渲染 el-pager>li
        if (page_no <= 5) {
            let p = ""
            for (let i = 1; i <= 5; i++) {
                if (page_no == i) {
                    p += "<li class='number active'>" + i + "</li>"
                } else {
                    p += "<li class='number'>" + i + "</li>"
                }
            }
            $(id + ">.el-pager").html(p)
        } else {
            let p = ""
            let now_max_page = Number($(id + ">.el-pager>li").eq(4).html())
            let now_min_page = Number($(id + ">.el-pager>li").eq(0).html())
            if (op == "+") {
                // 右移
                if (page_no > now_max_page) {
                    now_min_page += 1
                    now_max_page += 1
                }
            } else {
                // 左移
                if (page_no < now_min_page) {
                    now_min_page -= 1
                    now_max_page -= 1
                }
            }
            for (let i = now_min_page; i <= now_max_page; i++) {
                if (page_no == i) {
                    p += "<li class='number active'>" + i + "</li>"
                } else {
                    p += "<li class='number'>" + i + "</li>"
                }
            }
            $(id + ">.el-pager").html(p)
        }
    }
}
