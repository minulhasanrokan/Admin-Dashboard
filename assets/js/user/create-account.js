'use strict';

let signUpBtn = document.getElementById("sign_up");

if (signUpBtn) {
    signUpBtn.addEventListener("click", async function () {

        let controls = ["name#text#mandatory", "email#email#mandatory", "mobile#mobile#mandatory"];
        let msgTexts = ["Name", "Email Address", "Mobile Number"];

        if (document.getElementById("invite_code")) {
            controls.push("invite_code#text#mandatory");
            msgTexts.push("Invitation Code");
        }

        if (document.getElementById("password")) {
            controls.push("password#password#mandatory");
            controls.push("confirm_password#password#mandatory");
            msgTexts.push("Password");
            msgTexts.push("Confirm Password");
        }

        controls.push("terms_privacy#checkbox#mandatory");
        msgTexts.push("Terms & Privacy");

        let status = form_validation({
            control: controls.join("*"),
            msg_text: msgTexts.join("*")
        });

        if (status === false) {
            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        showLoader();

        let formData = new FormData();
        formData.append("name", $("#name").val());
        formData.append("email", $("#email").val());
        formData.append("mobile", $("#mobile").val());

        if (document.getElementById("invite_code")) {
            formData.append("invite_code", $("#invite_code").val());
        }

        if (document.getElementById("password")) {
            formData.append("password", $("#password").val());
            formData.append("confirm_password", $("#confirm_password").val());
        }

        formData.append("terms_privacy", $("#terms_privacy").is(":checked") ? "1" : "0");

        try {
            let response = await web.post('/user-registration/store', formData);

            if (response.status === 201) {
                showAlertBox({
                    title: "Success",
                    message: "User Registration Success",
                    type: "info",
                    sound: "/assets/sounds/success.mp3"
                });

                setTimeout(function () {
                    if (response.data && response.data.data && response.data.data.verified) {
                        window.location.href = BaseUrl + "/login";
                    } else {
                        window.location.href = BaseUrl + "/user-registration/verify-account/" + response.data.data.uuid;
                    }
                }, 3000);
            }
        } catch (err) {
            // web response interceptor displays server-side validation error
        }
    });
}

if (document.getElementById("name")) {
    document.getElementById("name").addEventListener("focusout", function () {
        form_validation({control: "name#text#mandatory", msg_text: "Name"});
    });
}

if (document.getElementById("email")) {
    document.getElementById("email").addEventListener("focusout", function () {
        form_validation({control: "email#email#mandatory", msg_text: "Email Address"});
    });
}

if (document.getElementById("mobile")) {
    document.getElementById("mobile").addEventListener("focusout", function () {
        form_validation({control: "mobile#mobile#mandatory", msg_text: "Mobile Number"});
    });
}

if (document.getElementById("invite_code")) {
    document.getElementById("invite_code").addEventListener("focusout", function () {
        form_validation({control: "invite_code#text#mandatory", msg_text: "Invitation Code"});
    });
}

if (document.getElementById("password")) {
    document.getElementById("password").addEventListener("focusout", function () {
        form_validation({control: "password#password#mandatory", msg_text: "Password"});
    });
}

if (document.getElementById("confirm_password")) {
    document.getElementById("confirm_password").addEventListener("focusout", function () {
        form_validation({control: "confirm_password#password#mandatory", msg_text: "Confirm Password"});
    });
}

if (document.getElementById("terms_privacy")) {
    document.getElementById("terms_privacy").addEventListener("click", function () {
        form_validation({control: "terms_privacy#checkbox#mandatory", msg_text: "Terms & Privacy"});
    });
}

if (document.getElementById("eye_btn_password")) {
    document.getElementById("eye_btn_password").addEventListener("click", function () {
        showHideValue('password', 'eye_btn_password');
    });
}

if (document.getElementById("eye_btn_c_password")) {
    document.getElementById("eye_btn_c_password").addEventListener("click", function () {
        showHideValue('confirm_password', 'eye_btn_c_password');
    });
}