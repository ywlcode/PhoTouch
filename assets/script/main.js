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
})

$()