$(function () {
    // 獲取即時狀態
    const url = '/dashboarda/datagrid';
    function getAPI() {
        $.get(url, (res) => {
            loading = false;
            var data = res.rows;
            // 跑狀態卡片
            data.forEach(item => {
                $('.board').append(`
                    <div class="card">
                        <div class="card-title card-header ${`index` + item.Id}" style="background:${statusColor(item.StatusColor)}">
                            <div class="number ${`index` + item.Id}">${item.MachineNumber}</div>
                            <div class="name ${`index` + item.Id}">${item.MachineName}</div>
                        </div>
                        <div class="card-content">
                            <div class="info info-head">
                                <div class="number ${`index` + item.Id}">${item.MoldNumber ? item.MoldNumber : ""}</div>
                                <div class="name ${`index` + item.Id}">${item.Cycletime ? item.Cycletime + 's': ""}</div>
                            </div>
                            <div class="progress-circle" >
                                <div class="circle ${`index` + item.Id}" style="border-color:${statusColor(item.StatusColor)}">
                                    <div class="${`progress` + item.Id}">${item.Progress ? item.Progress + '%' : ""}</div>
                                    <div class="${`planproduction` + item.Id}">${item.PlanProduction ? item.PlanProduction : ""}</div>
                                    <div class="${`plantimes` + item.Id}">${stringtotime(item.PlanTimes) ? stringtotime(item.PlanTimes) : ""}</div>
                                    <div class="${`status_text` + item.Id}">${showStatusText(item.MoldNumber, item.StatusColor)}</div>
                                </div>
                            </div>
                            <div class="info info-foot">
                                <div class="number ${`index` + item.Id}">${item.ExcuteQty}</div>
                                <div class="name ${`index` + item.Id}">${item.Qty}</div>
                            </div>
                        </div>
                        <div class="card-footer">
                            <div class="progress">
                                <div class="progress-bar ${`running` + item.Id}"" role="progressbar" style="width: ${stringtofloatconverse(item.Progress)}%; background:#3d91bc" aria-valuenow="${stringtofloatconverse(item.Progress)}" aria-valuemin="0" aria-valuemax="100"></div>
                            </div>
                        </div>
                    </div>
                `);
            });

            // 分頁模式
            var size = 24;
            var init = 0;
            var totalPage = Math.ceil(data.length / 24)
            console.log(totalPage)
            $('.Add').on('click', function () {
                init++;
                $('.pageBoard > div').remove();
                count(init, size, data, totalPage);
            })
            $('.reducer').on('click', function () {
                init--;
                $('.pageBoard > div').remove();
                if (init > -1) {
                    count(init, size, data, totalPage)
                } else {
                    init = 0;
                }
            })


            if (res) initInterface(true);
            count(init, size, data, totalPage);
            statistical(data)
        });

    };
    // 執行分頁
    function count(init, size, data, totalPage) {
        console.log(init)
        console.log(data)
        $('.reducer').show()
        $('.Add').show()
        if (init === 0) {
            $('.reducer').hide()
        }
        if (init === (totalPage - 1)) {
            $('.Add').hide()
        }
        $('.currnet').html(init + 1)
        $('.total').html(totalPage)
        var start = init * size;
        var end = start + size;
        // 使用陣列 slice 方法 取出資料
        data.slice(start, end).forEach(item => {
            $('.pageBoard').append(`
                    <div class="card">
                        <div class="card-title card-header ${`index` + item.Id}" style="background:${statusColor(item.StatusColor)}">
                            <div class="number ${`index` + item.Id}">${item.MachineNumber}</div>
                            <div class="name ${`index` + item.Id}">${item.MachineName}</div>
                        </div>
                        <div class="card-content">
                            <div class="info info-head">
                                <div class="number ${`index` + item.Id}">${item.MoldNumber ? item.MoldNumber : ""}</div>
                                <div class="name ${`index` + item.Id}">${item.Cycletime ? item.Cycletime + 's' : ""}</div>
                            </div>
                            <div class="progress-circle" >
                                <div class="circle ${`index` + item.Id}" style="border-color:${statusColor(item.StatusColor)}">
                                    <div class="${`progress` + item.Id}">${item.Progress ? item.Progress + '%' : ""}</div>
                                    <div class="${`planproduction` + item.Id}">${item.PlanProduction ? item.PlanProduction : ""}</div>
                                    <div class="${`plantimes` + item.Id}">${stringtotime(item.PlanTimes) ? stringtotime(item.PlanTimes) : ""}</div>
                                    <div class="${`status_text` + item.Id}">${showStatusText(item.MoldNumber, item.StatusColor)}</div>
                                </div>
                            </div>
                            <div class="info info-foot">
                                <div class="number ${`index` + item.Id}">${item.ExcuteQty}</div>
                                <div class="name ${`index` + item.Id}">${item.Qty}</div>
                            </div>
                        </div>
                        <div class="card-footer">
                            <div class="progress">
                                <div class="progress-bar ${`running` + item.Id}"" role="progressbar" style="width: ${stringtofloatconverse(item.Progress)}%; background:#3d91bc" aria-valuenow="${stringtofloatconverse(item.Progress)}" aria-valuemin="0" aria-valuemax="100"></div>
                            </div>
                        </div>
                    </div>
                `);
        })

    }

    // 更新狀態
    setInterval(function () {
        $.get(url, (res) => {
            var data = res.rows;
            data.forEach(item => {
                // card-header
                $(`.card > .card-title.index${item.Id}`).css({ 'background': statusColor(item.StatusColor) });
                $(`.card > .card-title > .number.index${item.Id}`).html(item.MachineNumber);
                $(`.card > .card-title > .name.index${item.Id}`).html(item.MachineName);
                // card-content
                // info-head
                $(`.card-content > .info-head > .number.index${item.Id}`).html(item.MoldNumber);
                $(`.card-content > .info-head > .name.index${item.Id}`).html(item.Cycletime + 's');
                // progress-circle
                $(`.card-content > .progress-circle > .circle.index${item.Id}`).css({ 'border-color': statusColor(item.StatusColor) });
                $(`.card-content > .progress-circle > .circle > .progress${item.Id}`).html(item.Progress ? item.Progress + "%" : " ");
                $(`.card-content > .progress-circle > .circle > .planproduction${item.Id}`).html(item.PlanProduction);
                $(`.card-content > .progress-circle > .circle > .plantimes${item.Id}`).html(stringtotime(item.PlanTimes));
                $(`.card-content > .progress-circle > .circle > .statustext${item.Id}`).html(showStatusText(item.MoldNumber, item.StatusColor));
                //  info-foot
                $(`.card-content > .info-foot > .number.index${item.Id}`).html(item.ExcuteQty);
                $(`.card-content > .info-foot > .name.index${item.Id}`).html(item.Qty);
                // card-footer
                $(`.card-footer > .progress > .holiday${item.Id}`).css({ 'width': item.ProgressHoliday + '%' });
                $(`.card-footer > .progress > .downtime${item.Id}`).css({ 'width': item.ProgressDowntime + '%' });
                $(`.card-footer > .progress > .idle${item.Id}`).css({ 'width': item.ProgressIdle + '%' });
                $(`.card-footer > .progress > .abnormal${item.Id}`).css({ 'width': item.ProgressAbnormal + '%' });
                $(`.card-footer > .progress > .running${item.Id}`).css({ 'width': stringtofloatconverse(item.Progress) + '%' });
            });

            statistical(data)
        });

    }, 10000);

    $('#fullscreen-head').hide();
    $('.show-fullscreen').hide();
    // 監聽視窗是否有改變大小
    window.addEventListener("resize", e => {
        // 如果 document 是全螢幕顯示 .text-full & #closefullscreen
        if (
            document.webkitFullscreenElement ||
            document.mozFullScreenElement ||
            document.msFullscreenElement ||
            document.fullscreenElement
        ) {
            $('#fullscreen-head').show();
            $('#fullscreen').addClass("fullscreen-bg-color");
            $('.show-fullscreen').show();
            $('.close-fullscreen').hide();
            $('#openfullscreen').hide();
            $('.statistical-board .small-box').addClass("addMargin")
        } else {
            // 如果 document 不是全螢幕隱藏 .text-full & #closefullscreen
            $('#fullscreen-head').hide();
            $('#fullscreen').removeClass("fullscreen-bg-color");
            $('.show-fullscreen').hide();
            $('.close-fullscreen').show();
            $('#openfullscreen').show();
            $('.statistical-board .small-box').removeClass("addMargin")
        }
    });


    initInterface();
    getAPI();
    currentTime();

    setInterval(function () {
        currentTime();
    }, 1000);

})

// 全螢幕監聽
$('#openfullscreen').on('click', (e) => {
    var elem = document.getElementById("fullscreen");

    $('.text-full').show();
    $('.show-fullscreen').show();
    $('.close-fullscreen').hide();
    if (elem.requestFullscreen) {
        elem.requestFullscreen();
    } else if (elem.msRequestFullscreen) {
        elem.msRequestFullscreen();
    } else if (elem.mozRequestFullScreen) {
        elem.mozRequestFullScreen();
    } else if (elem.webkitRequestFullscreen) {
        elem.webkitRequestFullscreen();
    }

});

// 統計狀態
function statistical(data) {
    // 統計卡片
    var total = data.length;
    var nurunning = 0;
    var idle = 0;
    var abnormal = 0;
    var produce = 0;
    data.forEach(item => {
        switch (item.StatusColor) {
            case "0":
                nurunning++
                break;
            case "1":
                nurunning++
                break;
            case "2":
                nurunning++
                break;
            case "3":
                idle++
                break;
            case "4":
                abnormal++
                break;
            case "5":
                produce++
                break;
            case "6":
                produce++
                break;
            default: return 0
        }
    })
    $('.text-total').html(total)
    $('.text-running').html(produce)
    $('.text-idle').html(idle)
    $('.text-abnormal').html(abnormal)
    $('.text-nowork').html(nurunning)

}


// loading 
function initInterface(loading) {
    if (loading) {
        $('#fullscreen').show();
        $('#screenLoading').hide();
    } else {
        $('#fullscreen').hide();
        $('#screenLoading').show();
    }
}

// 判斷狀態顏色
function statusColor(color) {
    var statusColor;
    switch (color) {
        case "0":
            statusColor = '#000000';
            break;
        case "1":
            statusColor = '#dcdcdc';
            break;
        case "2":
            statusColor = '#808080';
            break;
        case "3":
            statusColor = '#8b0000';
            break;
        case "4":
            statusColor = '#ff0000';
            break;
        case "5":
            statusColor = '#00a136';
            break;
        case "6":
            statusColor = '#0099ff';
            break;
        default:
            statusColor = "#808080";
    }
    return statusColor;
}

// 時間
function currentTime() {
    var currentTime = new Date();
    var year = currentTime.getFullYear();
    var month = currentTime.getMonth() + 1;
    var date = currentTime.getDate();
    var hour = currentTime.getHours();
    var min = currentTime.getMinutes();
    var sec = currentTime.getSeconds();

    timeString = zeroPrefixer(hour) + ":" + zeroPrefixer(min) + ":" + zeroPrefixer(sec);
    $('.current-time').html(timeString);
    // console.log(timeString);
}

// 判斷狀態文字
function showStatusText(schedule, status) {
    var str = false;
    if (!schedule) {
        switch (status) {
            case "0":
                str = '假期';
                break;
            case "1":
                str = '保養';
                break;
            case "2":
                str = '停機';
                break;
            case "3":
                str = "閒置"
                break;
            case "4":
                str = '異常';
                break;
            case "5":
                str = "運轉"
                break;
            case "6":
                str = '試模';
                break;
            default:
                str = " ";
        }
    }

    // console.log(schedule, status)

    return str ? str : " "
}

function zeroPrefixer(num) {
    return (num < 10 ? "0" : "") + num;
}

function stringtotime(str) {
    // var date = new Date(null);
    // var a = date.setSeconds(str);

    var timeint =parseInt(str)
    var day = Math.floor(timeint/86400)
    var hr = Math.floor((timeint%86400)/3600)
    var min = Math.floor(((timeint%86400)%3600)/60)
    var sec = ((timeint%86400)%3600)%60
    // var timeString = day + "d " + hr +"hr"

    if (day == 0){
        var timeString = hr +"hr " + min + "min";
    }
    else {
        var timeString = day + "d " + hr +"hr";
    }
    return timeString
    
}

function stringtofloatconverse(str) {
    var converseProgress = 100 - parseFloat(str);
    return converseProgress
}

// function stringtotime(str) {
//     var date = new Date(null);
//     date.setSeconds(str); // specify value for SECONDS here
//     var timeString = date.toISOString().substr(11, 8);
//     return timeString
// }