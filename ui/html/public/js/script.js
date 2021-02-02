function Logout() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/Logout', true);
    // r.setRequestHeader("Content-Type", "application/json");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        // let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                alert("로그아웃 되었습니다.");
                document.location.href = '/';
            } else {
                console.log("로그아웃 실패");
            }
        }
    };

    r.send();
}


// function CheckLogin() {
//     const r = new XMLHttpRequest();
//     r.open('GET', 'http://localhost:5000/api/CheckLogin', true);
//     r.withCredentials = true;
//     r.onreadystatechange = function () {
//         let responseObj;
//         if (r.readyState === 4) {
//             if (r.status === 200) {
//                 console.log("checkLogin 정상");
//                 responseObj = JSON.parse(r.responseText);
//                 $('#dropdown_id').text(responseObj.user_info.email);
//                 $('#dropdown_name').text(responseObj.user_info.name);
//                 document.location.href = "./dashboard/main.html"
//             } else {
//                 document.location.href = '/';
//             }
//         }
//
//     };
//     r.send();
// }

function CheckLoginInLoginPage() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/checklogin', true);
    r.withCredentials = true;
    r.onreadystatechange = function () {
        // let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                console.log("checkLogin 정상");
                // responseObj = JSON.parse(r.responseText);
                // $('#dropdown_id').text(responseObj.user_info.email);
                // $('#dropdown_name').text(responseObj.user_info.name);
                document.location.href = "./dashboard/main.html"
            } else {
                document.location.href = '/';
            }
        }

    };
    r.send();
}

function Tg_total() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/setting/getTag', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        let rObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                rObj = JSON.parse(r.responseText);
                if (rObj.tags.length > 0) {
                    for (let j = 0; j < rObj.tags.length; j++) {
                        tg_total[j] = rObj.tags[j].tag_name;
                    }
                }
            } else if (r.status === 403) {
                const tokenResult = Refresh();
                if (tokenResult) {
                    Tg_total();
                    console.log("true")
                }
            } else {
                document.location.href = '/';
            }
        }
    };
    r.send();
}

function Login() {
    const r = new XMLHttpRequest();
    r.open('POST', 'http://localhost:5000/api/Login', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        let rObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                document.location.href = './dashboard/total'
            } else if (r.status === 400) {
                alert("계정 정보를 입력해주세요. ");
            } else if (r.status === 401) {
                alert("계정 정보를 확인해주세요. ");
            } else if (r.status === 402) {
                alert("이메일, 비밀번호 형식이 잘못됐습니다. ");
            } else if (r.status === 403) {
                alert("존재하지 않는 계정입니다.");
            } else if (r.status === 500) {
                alert("서버에러");
            } else if (r.status === 408) {
                alert("로그인 실패 횟수를 초과했습니다(5회)\n관리자에게 연락해주세요.")
            }
        }
    };
    const formData = $('form').serializeObject();
    r.send(JSON.stringify(formData));
}

function GetDashBoard() {
    const user_id = document.getElementById("dropdown_id");
    const user_name = document.getElementById("dropdown_name");
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/dashboard', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                console.log(r.responseText);
                responseObj = JSON.parse(r.responseText);
                user_id.innerHTML = responseObj.email
                user_name.innerHTML = responseObj.name
                // user_id.text(responseObj.email);
                // user_name.text(responseObj.name);
            } else if (r.status === 403) {
                const tokenResult = Refresh();
                if (tokenResult) {
                    GetDashBoard()
                    console.log("true")
                }

            } else {
                document.location.href = '/';
            }
        }
    };

    r.send();
}

function Register() {
    if (!$("input:checkbox[name='checkbox']").is(":checked")) {
        alert("이용약관에 동의해주세요");
        return;
    }
    const r = new XMLHttpRequest();
    r.open('POST', 'http://localhost:5000/api/createUser', true);
    r.onreadystatechange = function () {
        if (r.readyState === 4) {
            if (r.status === 200) {
                alert("회원가입이 완료되었습니다.");
                document.location.href = '/'
            } else if (r.status === 400) {
                alert("회원정보를 입력해주세요. ");
            } else if (r.status === 401) {
                alert("비밀번호를 확인해 주세요.");
            } else if (r.status === 402) {
                alert("이메일, 비밀번호 형식을 확인해 주세요. ");
            } else if (r.status === 403) {
                alert("이미 존재하는 계정입니다. ");
            } else if (r.status === 500) {
                alert("서버에러")
            }

        }
    };
    const formData = $('#register_form').serializeObject();
    r.send(JSON.stringify(formData));
}

function getRoot() {
    const r = new XMLHttpRequest();
    r.open('GET', '/', true);
    r.send()
}

function Refresh() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/RefreshToken', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        if (r.readyState === 4) {
            if (r.status === 200) {
                console.log(r.responseText);
                return true;
            } else {
                alert("인증 토큰 갱신 실패했습니다.")
                document.location.href = "/";
            }
        }
    };
    const formData = $('form').serializeObject();
    r.send(JSON.stringify(formData));
}
