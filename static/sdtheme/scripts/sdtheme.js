/*!
 * 用於當前主題通過的方法
 * Created by lihaitao on 2017-7-10.
 */
layer.ready(function () {
    layer.config({
        isOutAnim: false,
       /* extend: 'metronic/style.css',
        skin: 'layer-ext-metronic'*/
    });
});

/**
 * Created by lihaitao on 2017-7-11.
 */
var sdtheme = function () {
    showstr = function (str, replace) {
        if (str === null || typeof (str) === "undefined") {
            if (typeof (replace) === 'undefined') {
                return "";
            } else {
                return replace;
            }

        }
        return str;
    }
    showlongstr = function (str, replace) {
        return '<span title="' + str + '">' + showstr(str, replace) + '</span>';
    }
    showenable = function (val) {
        if (val === 1 || val === "1") {
            return '<label class="label label-success label-sm"><i class="fa fa-check"></i> 啟用</label>';
        } else if (val === 0 || val === "0") {
            return '<label class="label label-danger label-sm"><i class="fa fa-ban"></i> 禁用</label>';
        } else if (val === -1 || val === "-1")
            return '<label class="label label-info label-sm"><i class="fa fa-trash"></i> 刪除</label>';
        else {
            return "";
        }
    }
    showyes = function (val) {
        if (val === 1 || val === "1" || val === true) {
            return '<label class="label label-primary label-sm"><i class="fa fa-check"></i> 是</label>';
        } else if (val === 0 || val === "0" || val === false) {
            return '<label class="label label-danger label-sm"><i class="fa fa-close"></i> 否</label>';
        } else {
            return "";
        }
    }
    showenum = function (value, texts, css, icon) {
        var index = 0, text = "", icss = 'label-default';
        if (css === null || typeof (css) === 'undefined') {
            css = ['label-primary', 'label-success', 'label-info', 'label-warning', 'label-danger', 'label-default'];
        }
        var item = $(texts).each(function (i, e) {
            if (e.Value === value) {
                index = i;
                text = e.Text;
                return false;
            }
        });
        if (index <= css.length) {
            icss = css[index];
        }
        return '<label class="label ' + icss + '  label-sm">' + text + '</label>';
    }

    //從cookie加載查詢條件
    //默認存入cookie的key由formId而，可以使用extKey進一步區分
    //callback用於，加載完數據後，做更多事件，比如處理級聯的下拉框
    function loadSearchText(formId, extKey, callback) {
        if(!extKey){
            extKey = ''
        }
        var serialize = $.cookie('formmaitain_' + formId + extKey);
        if (serialize) {
            serialize = serialize.replace(/\+/g, ' ');
            //整理name 和 值，
            var namevals = {};
            $(serialize.split('&')).each(function (i, e) {
                var keyval = e.split('=');
                if (namevals[keyval[0]] !== undefined) {
                    namevals[keyval[0]] = namevals[keyval[0]] + ',' + keyval[1]; //考慮同一個name多個值的情況
                } else {
                    namevals[keyval[0]] = keyval[1];
                }
            });
            for (var key in namevals) {
                var ctrl = $("#" + formId).find('[name="' + key + '"]');
                if (ctrl.length > 0) {
                    ctrl = ctrl.get(0);
                    if (ctrl.tagName.toLowerCase() == "input") {
                        if (($(ctrl).prop('type') === 'checkbox' || $(ctrl).prop('type') === 'radio') && !$(ctrl).prop('checked')) {
                            $(ctrl).parent().trigger('click');
                        } else {
                            $(ctrl).val(namevals[key]);
                        }
                    } else if (ctrl.tagName.toLowerCase() == "select") {
                        $(ctrl).selectpicker('val', namevals[key].split(',')); //將,拼接的字符轉成數組
                    }
                }
            }
            if(callback){
                callback(namevals)
            }
        }
    }
    //將查詢條件存入cookie
    function saveSearchText(formId , extKey) {
        if(!extKey){
            extKey = ''
        }
        //將查詢表單的值存在cookie
        $.cookie('formmaitain_' + formId + extKey, decodeURIComponent($("#" + formId).serialize(), true), { expires: 1 });
    }
     //treetable
     function saveExpandStatus(treeGridId, extKey) {
        if(!extKey){
            extKey = ''
        }
        //獲取所有展開的節點的id
        var ids = [];
        var $expandeds = $("#" + treeGridId).find('tr.expanded');
        $expandeds.each(function (i, e) {
            ids.push($(e).attr('data-tt-id'));
        });
        $.cookie(treeGridId + '_expandedids' + extKey, ids.join(','), { expires: 1 });
    }
    //treetable
    function loadExpandStatus(treeGridId, extKey) {
        if(!extKey){
            extKey = ''
        }
        //獲取所有展開的節點的id
        var ids = $.cookie(treeGridId + '_expandedids' + extKey);
        if (ids !== null && typeof ids !== 'undefined' && ids.length > 0) {
            $(ids.split(',')).each(function (i, e) {
                //先判斷節點是否存在
                if ($("#" + treeGridId).find('tr[data-tt-id="' + e + '"]').length > 0) {
                    $("#" + treeGridId).treetable('expandNode', e);
                }
            });
        }
    }
    //高亮顯示，scrollto是否自動滾動， 傳入jquery對象,css
    function highlight(object, scrollto, css) {
        if (object === null || object.length === 0) return;
        var t = 6, scroll = true, hcss = 'highlight';
        if (scrollto !== null && typeof scrollto !== 'undefined') {
            scroll = scrollto
        }
        if (css !== null && typeof css !== 'undefined') {
            hcss = css
        }
        //滾動條自動滾動
        var $win = $(window);
        var $body = $('html,body');
        var troffsettop = object.offset().top;
        if (troffsettop < $win.scrollTop() + object.outerHeight() * 2) {
            $body.stop().animate({ "scrollTop": troffsettop - object.outerHeight() * 2 }, 200);
        }
        if (troffsettop >= $win.scrollTop() + $win.height() - object.outerHeight() * 3) {
            $body.stop().animate({ "scrollTop": troffsettop - $win.height() + object.outerHeight() * 3 }, 200);
        }
        //高亮
        $(object).toggleClass(hcss);
        var spark = function () {
            if (t-- === 0) {
                $(object).removeClass(hcss);
                return;
            }
            $(object).toggleClass(hcss);
            setTimeout(spark, 300);
        }
        spark();
    }
    //初始化搜索條件面板狀態保持功能
    function searchPanelStatusInit(btnid, hidetips, extKey) {
        if(!extKey){
            extKey = ''
        }
        var $btn = $('#' + btnid);
        if ($btn.length > 0) {
            var $icon = $('i',$btn)
            //在點擊事件裡保存狀態到cookie
            $btn.off('click').on('click', function () {
                console.log('click')
                //點擊時保存，css會切換
                var css = 'fa-plus';
                if ($icon.hasClass('fa-plus')) {
                    css = 'fa-minus';
                }
                $.cookie('SearchPanelStatus' + btnid + extKey, css, { expires: 1 });
            });
            //頁面加載時，加載狀態
            var css = $.cookie('SearchPanelStatus' + btnid + extKey);
            if (css != null && typeof css !== 'undefined') {
                if (css === 'fa-minus') {
                    $icon.removeClass('fa-plus').addClass('fa-minus')
                    $btn.closest('div.box').removeClass('collapsed-box')
                    $btn.closest('div.box-header').next().show()
                }
            }
            if(!hidetips || hidetips===false){
                //只要面板處於關閉
                if($icon.hasClass('fa-plus')){
                    //重點提示更多條件
                    layer.tips('顯示/隱藏更多查詢條件', '#' + btnid, {
                        tips: [1, '#35AA47'],
                        time: 5000,
                        shift: 6
                    });
                }
            }
        }
    }
    //將當前滾動條位置保存至cookie,默認會話結束後失效
    function saveScrollTop2Cookie(expire) {
        var exp = null;
        if (expire !== null && typeof expire !== 'undefined') {
            exp = expire;
        }
        var scrollTop = $(window).scrollTop();
        if (scrollTop > 0) {
            $.cookie('page.scrollTop', scrollTop, { expires: exp });
        }
    }
    //從cookie讀取滾動條位置，使用一次後失效
    function loadScrollTopFromCookie() {
        var scrollTop = $.cookie('page.scrollTop');
        if (scrollTop != null && typeof scrollTop !== 'undefined' && scrollTop.length > 0 && parseInt(scrollTop) > 0) {
            $(window).scrollTop(parseInt(scrollTop));
            $.cookie('page.scrollTop', '');
        }
    }
    function alertXHRError (XHR, status, e) {
        if (typeof (XHR.responseText) !== 'undefined') {
            parent.layer.alert("請求失敗：" + XHR.responseText, { icon: 2, title: '錯誤' });
        }
        else {
            parent.layer.alert("請求失敗", { icon: 2, title: '錯誤' });
        }
    }
    //開始和結束時間級聯
    function timeCtrlBeginToEnd(beginselector,endselector,ctrtype,format){
        if(!ctrtype) {
            ctrtype = 'datetime'
        }
        if (!format) {
            format = 'YYYY/MM/DD HH:mm:ss'
        }
        //時間範圍控制
        var dateFirst = $(beginselector);
        var dateEnd = $(endselector);
        var dateFirstApi = null;
        var dateEndApi = null;
        var EndTimeInit = function () {
            //獲取開始時間
            var startTime = parseInt(dateFirstApi.getDate('TIME'), 10);
            //獲取開始時間的日期
            var startTimeDate = dateFirstApi.getDate('YYYY/MM/DD');
            var endTime = parseInt(dateEndApi.getDate('TIME'), 10);
            if (endTime < startTime) {
                dateEndApi.clearDate();
            };
            dateEndApi.setOptions({
                //startDate: startTime //這種方式會導致日期加1
                startDate: startTimeDate
            });
            dateEndApi.show();
        }
        dateFirst.bind('change', function () {
            EndTimeInit();
        }).cxCalendar(function (api) {
            dateFirstApi = api;
            dateFirstApi.setOptions({
                type: ctrtype,
                format: format
            });
        });
        dateEnd.bind('click', function () {
            EndTimeInit();
        }).cxCalendar(function (api) {
            dateEndApi = api;
            dateEndApi.setOptions({
                type: ctrtype,
                format: format
            });
        });
    }
    return {
        //
        init:init,
        //
        uniform:uniform,
        //傳入的值為null，則用replace代替，默認為空
        showstr: showstr,
        //使用span將值包裹
        showlongstr: showlongstr,
        //顯示啟用或者禁用
        showenable: showenable,
        //顯示是否
        showyes: showyes,
        //顯示枚舉
        showenum: showenum,
        //保存form裡的查詢條件
        saveSearchText: saveSearchText,
        //加載form查詢條件
        loadSearchText: loadSearchText,
        //將treetable裡展開的節點數據保存到cookie
        saveExpandStatus: saveExpandStatus,
        //從cookie讀取treetable展開節點數據，並應用
        loadExpandStatus: loadExpandStatus,
        //高亮顯示
        highlight: highlight,
        //初始化搜索條件面板狀態保持功能
        searchPanelStatusInit: searchPanelStatusInit,
        //將滾動條位置存入cookie
        saveScrollTop2Cookie: saveScrollTop2Cookie,
        //加載滾動條位置
        loadScrollTopFromCookie: loadScrollTopFromCookie,
        //提示錯誤
        alertXHRError:alertXHRError,
        //時間區間初始化
        timeCtrlBeginToEnd:timeCtrlBeginToEnd
    };
    //控件美化
    function uniform() {
        //使用bootstrap-select的下拉初始化
        $("select.bs-select").selectpicker();
    }
    function init() {
        //控件美化
        uniform();
    }
}();

jQuery(document).ready(function () {
    sdtheme.init()
});