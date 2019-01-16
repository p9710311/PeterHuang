var rms = function () {
    //初始化    
    function init() {
        //console.log("*** rms js *** init");
    }
    //菜單初化
    function pageSidebarInit(options) {
        var url = options.url;
        $.sdpost(url, {}, function (re) {
            if (re.code === 0) {                
                var $pageSidebar = $(options.slideBarBox);
                if ($pageSidebar.length === 0) {
                    console.log("menu box is null");
                    return;
                }
                $pageSidebar.html('');
                var data = re.obj;
                var html = [];
                $(data).filter(function (i, e) {
                    //找出要節點
                    return e.Parent.Id === 0;
                }).each(function (i, e) {
                    //菜單
                    if (e.Rtype === 1) {
                        //遞歸加載子節點
                        html.push(showSiderBarSon(e, data));
                    } else if (e.Rtype === 0) {
                        //如果是區域，先添加header
                        html.push('<li class="header" > ' + e.Name + ' </li>');
                        //添加區域的子節點  
                        $(data).filter(function (i1, e1) {
                            return e1.Parent.Id === e.Id;
                        }).each(function (i2, e2) {
                            //遞歸調用，顯示子節點
                            html.push(showSiderBarSon(e2, data));
                        });
                    }
                });
                $pageSidebar.html(html.join(''));
                //acitve 將href值與cur對應上的菜單激活
                //console.log(options.cur); 
                var curli = $('li a[href="'+options.cur+'"]',$pageSidebar);
                if(curli.length>0){
                    //激活當前
                    curli.parent().addClass('active');  
                    //所有父ul.treeview-menu顯示                
                    curli.parents("ul.treeview-menu").show();
                    //所有父li.treeview展開
                    curli.parents('li.treeview').addClass('menu-open');     
                }      
            } else {
                layer.alert("獲取失敗", { icon: 2, title: "失敗" })
            }
        });
    }    
    function showSiderBarSon(cur, data) {
        var html = [];
        //有子菜單
        if (cur.SonNum > 0) {
            html.push('<li class="treeview">');
            html.push('  <a href="javascript:;">');
            html.push('    <i class="' + getIcon(cur.Icon) + '"></i>');
            html.push('    <span>' + cur.Name + '</span>');
            html.push('    <span class="pull-right-container"><i class="fa fa-angle-left pull-right"></i></span>');
            html.push('  </a>');
            html.push('  <ul class="treeview-menu">');
            $(data).filter(function (i, e) {
                return e.Parent.Id === cur.Id;
            }).each(function (i, e) {
                //遞歸調用，顯示子節點
                html.push(showSiderBarSon(e, data));
            });
            html.push('  </ul>')
            html.push('</li>')
        } else {
            if (!cur.LinkUrl || cur.LinkUrl.length === 0) {
                cur.LinkUrl = "javascript:;"
            }
            //無子菜單
            html.push('<li><a href="' + cur.LinkUrl + '"><i class="' + getIcon(cur.Icon) + '"></i><span>' + cur.Name + '</span></a></li>');
        }
        return html.join('');
    }
    function getIcon(icon) {
        if (!icon || icon.length === 0) {
            return "fa fa-circle-o"
        }
        return icon;
    }
    //全選和單選聯動
    chkAllSingleInit = function (chkAllSelector, chkSingleSelector) {
        $(chkAllSelector).on('click', function () {
            $(chkSingleSelector + ':enabled').prop('checked', $(this).prop('checked'));
        });
        $(chkSingleSelector).on('click', function () {
            if ($(this).prop('checked') === false) {
                $(chkAllSelector).prop('checked', false);
            } else {
                var totalnum = $(chkSingleSelector + ':enabled').length;
                var checknum = $(chkSingleSelector + ":enabled:checked").length;
                $(chkAllSelector).prop('checked', totalnum === checknum);
            }
        });
        //第一次加載時判斷是否全選了
        var totalnum = $(chkSingleSelector + ':enabled').length;
        var checknum = $(chkSingleSelector + ":enabled:checked").length;
        $(chkAllSelector).prop('checked', totalnum === checknum);
    }
    //獲取x-editable所需要的參數表,url為更新值時請求服務器地址，validate驗證的方式，refreshPk成功後是否根據主鍵刷新列表
    getEditableParam = function (url, validate, refreshPk) {
        console.log('validate:'+validate);
        console.log('refreshPk'+refreshPk);
        var defaultvalidate = function (value) {
            if ($.trim(value).length === 0) {
                return "必填項";
            }
            if (!/^\d+$/.test(value)) {
                return "必須為正整數，且不大於999999";
            }
        };
        if (validate !== null && typeof (validate) !== 'undefined') {
            defaultvalidate = validate
        }
        return {
            url: url,
            type: 'text',
            onblur: 'cancel', //Action when user clicks outside the container. Can be cancel|submit|ignore.  Setting ignore allows to have several containers open.
            showbuttons: true,
            ajaxOptions: {
                type: 'post',
                dataType: 'json'
            },
            validate: defaultvalidate,
            success: function (response, newValue) {
                if (response.code == 0) {
                    parent.layer.msg(response.msg, { icon: 1, title: '成功' });                   
                    if (refreshPk) {
                        //刷新列表數據
                        refresh(response.obj);
                    }
                    else {
                        refresh();
                    }                    
                } else {
                    return response.msg;
                }
            },
            error: function (response, newValue) {
                if (response.status === 500) {
                    return '系統錯誤，請刷新後重試';
                } if (response.status === 404) {
                    return '請求失敗，請刷新後重試';
                } else {
                    return response.responseText;
                }
            }
        }
    }
    return {
        init: init,
        //頁面左側菜單初始化
        pageSidebarInit: pageSidebarInit,
        //全選 單選初始化
        chkAllSingleInit:chkAllSingleInit,
        //獲取Editable插件的參數
        getEditableParam:getEditableParam
    }

}();
jQuery(document).ready(function () {
    rms.init()
});