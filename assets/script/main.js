var kk = 200.0;
var mm = 80.0;


$(window).scroll(function(){
    if(window.scrollY > kk) {
        var k = document.querySelector(".show").cloneNode(true);
        var index = document.querySelector(".indeximgs");
        index.appendChild(k);
        index.appendChild(k);
        index.appendChild(k);
        index.appendChild(k);
        kk += mm;
        console.log(window.scrollY);
    };
});

var flag = 1;

$(document).ready(function(){
    $("#z2").click(function(){
        if(flag == 1)
        {
            $(".content").css("width","calc(100% - 180px)");
            $(".content").css("left","180px");
            $(".side").css("width","180px");
            $(".box span").css("opacity","1");
            $("#z2 i").addClass("icon-fanhuijiantou");
            $("#z2 i").removeClass("icon-chufadaodaxiao");
            flag = 0;
        }
        else {
            $(".content").css("width","calc(100% - 78px)");
            $(".content").css("left","78px");
            $(".side").css("width","78px");
            $(".box span").css("opacity","0");
            $("#z2 i").addClass("icon-chufadaodaxiao");
            $("#z2 i").removeClass("icon-fanhuijiantou");
            flag = 1;
        }
    })
});