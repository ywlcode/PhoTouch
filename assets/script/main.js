var kk = 200.0;
var mm = 80.0;


$(window).scroll(function () {
    if (window.scrollY > kk) {
        var index = document.querySelector(".indeximgs");
        var k = document.querySelector(".show").cloneNode(true);
        $.getJSON("/img/rand", function (data) {
            for (var i = 0; i < data.length; i++) {
                index.appendChild(k);
                index.lastElementChild.firstElementChild.src = data[i].minurl;
            }
        })
        kk += mm;
    };
});

var flag = 1;

$(document).ready(function () {
    $("#z2").click(function () {
        if (! $.cookie('user')) {
            $("#userkk").text("未登录");
        }
        else {
            $("#userkk").text("您好");
        }
        if (flag == 1) {
            $(".content").css("width", "calc(100% - 180px)");
            $(".content").css("left", "180px");
            $(".side").css("width", "180px");
            $(".box span").css("opacity", "1");
            $("#z2 i").addClass("icon-fanhuijiantou");
            $("#z2 i").removeClass("icon-chufadaodaxiao");
            flag = 0;
        }
        else {
            $(".content").css("width", "calc(100% - 78px)");
            $(".content").css("left", "78px");
            $(".side").css("width", "78px");
            $(".box span").css("opacity", "0");
            $("#z2 i").addClass("icon-chufadaodaxiao");
            $("#z2 i").removeClass("icon-fanhuijiantou");
            flag = 1;
        }
    })
});


$(document).ready(function () {
    $("body").on("click", ".show img", function () {
        //获取页面高度和宽度
        var iWidth = document.documentElement.clientWidth;
        var iHeight = document.documentElement.clientHeight;
        //创建背景层
        var bgObj = document.createElement("div");
        bgObj.setAttribute("id", "bgbox");
        bgObj.style.width = iWidth + "px";
        bgObj.style.height = Math.max(document.body.clientHeight, iHeight) + "px";
        document.querySelector(".main").appendChild(bgObj);
        var oShow = document.getElementById('tanchu');
        var uu = $(this).attr("src"); //.children("img")
        $.post("/img/big", { minurl: uu }, function (data) {
            $("#tanchu img").attr("src", data)
        });
        document.body.style.overflowY = 'hidden';
        oShow.style.display = 'block';
        oShow.style.width = iWidth + "px";
        oShow.style.height = iHeight + "px";

        function oClose() {
            oShow.style.display = 'none';
            document.body.style.removeProperty('overflow-y');
            document.querySelector(".main").removeChild(bgObj);
            $("#tanchu img").attr("src", "")
        }

        var oClosebtn = document.createElement("span");
        oClosebtn.innerHTML = "×";
        oClosebtn.style.fontSize = "65px";
        oClosebtn.style.color = "white"
        oShow.appendChild(oClosebtn);
        oClosebtn.onclick = oClose;
        //bgObj.onclick = oClose;
    })
})

$(document).ready(function () {
    $("body").on("click", ".upload", function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        //获取页面高度和宽度
        var iWidth = document.documentElement.clientWidth;
        var iHeight = document.documentElement.clientHeight;
        //创建背景层
        var bgObj = document.createElement("div");
        bgObj.setAttribute("id", "bgbox");
        bgObj.style.width = iWidth + "px";
        bgObj.style.height = Math.max(document.body.clientHeight, iHeight) + "px";
        document.querySelector(".main").appendChild(bgObj);
        var oShow = document.querySelector(".upup")
        document.body.style.overflowY = 'hidden';
        oShow.style.display = 'block';
        function oClose() {
            oShow.style.display = 'none';
            document.body.style.removeProperty('overflow-y');
            document.querySelector(".main").removeChild(bgObj);
        }
        var oClosebtn = document.querySelector("#closebtn");
        oClosebtn.style.fontSize = "65px";
        oClosebtn.style.color = "white"
        oClosebtn.onclick = oClose;
    })
})

$(document).ready(function () {
    $(".blackbox i").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        $(this).toggleClass("icon-xihuan1");
        $(this).toggleClass("icon-xihuan");
        var uu = $(this).parents(".show").children("img").attr("src");
        $.post("/good", { minurl: uu });
    })
})

var lastpp = $("#z3");

$(document).ready(function () {
    $("#z3").css("border","1px solid white");
    lastpp = $("#z3");
    $("#z3").click(function () {
        window.location.href = '/';
    })
});

$(document).ready(function () {
    $("#z4").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});

$(document).ready(function () {
    $("#z5").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});

$(document).ready(function () {
    $("#z6").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});

$(document).ready(function () {
    $("#z7").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});

$(document).ready(function () {
    $("#z8").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});

$(document).ready(function () {
    $("#z9").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});

$(document).ready(function () {
    $("#z10").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});

$(document).ready(function () {
    $("#z11").click(function () {
        if (! $.cookie('user')) {
            alert("请登录");
            var t = setTimeout(function () {
                window.location.href = '/login';
            }, 300);
            return;
        }
        $.removeCookie('user', { path: '/' });
        $("#userkk").text("未登录");
        lastpp.css("border","");
        $(this).css("border","1px solid white");
        lastpp = $(this);
    })
});