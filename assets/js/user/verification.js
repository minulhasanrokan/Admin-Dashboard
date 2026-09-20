document.getElementById("verify_account").addEventListener("click", async function () {

    let status = form_validation({
        control: "digit_one#number#mandatory*digit_two#number#mandatory*digit_three#number#mandatory*digit_four#number#mandatory",
        msg_text: "Digit One*Digit Two*Digit Three*Digit Four"
    });

    if (status === false) {

        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});

        return false;
    }

    showLoader();

    let formData = new FormData();

    formData.append("uuid", $("#uuid").val());
    formData.append("digit_one", $("#digit_one").val());
    formData.append("digit_two", $("#digit_two").val());
    formData.append("digit_three", $("#digit_three").val());
    formData.append("digit_four", $("#digit_four").val());

    let response = await web.post('/user-registration/verify-store', formData);

    if (response.status === 200) {

        showAlertBox({
            title: "Success",
            message: "User Verification Success",
            type: "info",
            sound: "/assets/sounds/success.mp3"
        });

        setTimeout(function () {
            window.location.href = BaseUrl;
        }, 5000);
    }
});

let timerEl = document.getElementById("timer");
let resendBtn = document.getElementById("resend_otp");

let minutes = parseInt(timerEl.getAttribute("data-minutes"));
let time = minutes * 60;

resendBtn.classList.add("disabled-link");

let interval = startTimer();

resendBtn.addEventListener("click", async function (e) {
    e.preventDefault();

    if (this.classList.contains("disabled-link")) return;

    let formData = new FormData();

    formData.append("uuid", $("#uuid").val());

    let response = await web.post('/user-registration/resend-token', formData);

    if (response.status === 200) {

        showAlertBox({
            title: "Success",
            message: "Resend OTP Success",
            type: "info",
            sound: "/assets/sounds/success.mp3"
        });

        time = minutes * 60;
        this.classList.add("disabled-link");

        interval = startTimer();
    }
});