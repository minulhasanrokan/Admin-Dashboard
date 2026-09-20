'use strict';

function initProfilePage() {
    bindProfileEvents();
}

if (document.readyState === 'loading') {
    $(document).ready(initProfilePage);
} else {
    initProfilePage();
}

function bindProfileEvents() {

    let form = $("#form-edit-profile");
    let passwordForm = $("#form-change-password");

    form.off(".profileSettings");
    passwordForm.off(".passwordSettings");
    $("#profile-wrapper").off(".profileContainer");
    $(document).off(".profileTabSwitch");

    // Input validation triggers on focusout (matching business-settings.js)
    form.on("focusout.profileSettings", "input, textarea", function () {
        validateField(this.id);
    });

    passwordForm.on("focusout.passwordSettings", "input", function () {
        validatePasswordField(this.id);
    });

    // Form submission buttons
    form.on("click.profileSettings", "#profile_setting_button", async function () {
        await saveProfileSettings();
    });

    passwordForm.on("click.passwordSettings", "#password_change_button", async function () {
        await savePasswordSettings();
    });

    // Tab Navigation - Top tabs, In-Page "Edit" button, and "Cancel" buttons
    $(document).on("click.profileTabSwitch",
        "#profile-top-tabs .nav-link, .profile-nav-trigger, .tab-switcher-btn",
        function (e) {
            e.preventDefault();
            let targetPane = $(this).data("target-pane") || $(this).data("bs-target") || $(this).attr("href");
            if (targetPane) {
                switchProfileTab(targetPane);
            }
        }
    );

    // Password visibility toggle
    $("#profile-wrapper").on("click.profileContainer", ".toggle-password-btn", function (e) {
        e.preventDefault();
        let inputId = $(this).data("input-id");
        let input = $("#" + inputId);
        let icon = $(this).find("i");

        if (input.attr("type") === "password") {
            input.attr("type", "text");
            icon.removeClass("ti-eye").addClass("ti-eye-off");
        } else {
            input.attr("type", "password");
            icon.removeClass("ti-eye-off").addClass("ti-eye");
        }
    });

    // Photo upload trigger
    $("#profile-wrapper").on("click.profileContainer", "#btn-trigger-photo-upload", function (e) {
        e.preventDefault();
        $("#profile_photo_input").trigger("click");
    });

    $("#profile-wrapper").on("change.profileContainer", "#profile_photo_input", async function () {
        await uploadProfilePhoto(this);
    });
}

function validateField(id) {

    switch (id) {
        case "name":
            form_validation({control: "name#text#mandatory", msg_text: "Name"});
            break;
        case "mobile":
            form_validation({control: "mobile#mobile#nonmandatory", msg_text: "Mobile Number"});
            break;
        case "address":
            form_validation({control: "address#text#nonmandatory", msg_text: "Address"});
            break;
        case "details":
            form_validation({control: "details#text#nonmandatory", msg_text: "About"});
            break;
    }
}

function validatePasswordField(id) {

    switch (id) {
        case "current_password":
            form_validation({control: "current_password#text#mandatory", msg_text: "Current Password"});
            break;
        case "password":
            form_validation({control: "password#password#mandatory", msg_text: "New Password"});
            break;
        case "confirm_password":
            form_validation({control: "confirm_password#password#mandatory", msg_text: "Confirm Password"});
            break;
    }
}

async function saveProfileSettings() {

    let status = form_validation({
        control: "name#text#mandatory*mobile#mobile#nonmandatory*address#text#nonmandatory*details#text#nonmandatory",
        msg_text: "Name*Mobile Number*Address*About"
    });

    if (status === false) {
        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
        return false;
    }

    let confirmed = await confirmDialog("Profile Settings", "Are you sure you want to save the Profile Settings?");

    if (!confirmed) {
        return;
    }

    let formData = new FormData();
    formData.append("name", $("#name").val());
    formData.append("mobile", $("#mobile").val());
    formData.append("address", $("#address").val());
    formData.append("details", $("#details").val());

    try {
        let response = await web.post("/profile/update", formData);

        if (response.status === 200) {

            showAlertBox({title: "Success", message: "Profile Settings Updated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});

            let name = $("#name").val();
            let mobile = $("#mobile").val();
            let address = $("#address").val();
            let details = $("#details").val();

            $("#profile-card-name").text(name);
            $("#profile-detail-name").text(name);
            $("#profile-detail-mobile").text(mobile || "-");
            $("#profile-detail-address").text(address || "No address provided");
            $("#profile-detail-details").text(details || "No bio provided");
            $(".profile-dropdown h5, .header-user h6").text(name);

            switchProfileTab("pane-overview");
        }
    } catch (err) {
        // web response interceptor displays server-side validation error
    }
}

async function savePasswordSettings() {

    let status = form_validation({
        control: "current_password#text#mandatory*password#password#mandatory*confirm_password#password#mandatory",
        msg_text: "Current Password*New Password*Confirm Password"
    });

    if (status === false) {
        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
        return false;
    }

    let confirmed = await confirmDialog("Change Password", "Are you sure you want to update your password?");

    if (!confirmed) {
        return;
    }

    let formData = new FormData();
    formData.append("current_password", $("#current_password").val());
    formData.append("new_password", $("#password").val());
    formData.append("confirm_password", $("#confirm_password").val());

    try {
        let response = await web.post("/profile/change-password", formData);

        if (response.status === 200) {

            showAlertBox({title: "Success", message: "Password Updated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});

            $("#current_password").val("");
            $("#password").val("");
            $("#confirm_password").val("");

            switchProfileTab("pane-overview");
        }
    } catch (err) {
        // web response interceptor displays server-side validation error
    }
}

async function uploadProfilePhoto(input) {

    let file = input.files[0];
    if (!file) return;

    let allowedTypes = ["image/jpeg", "image/png", "image/webp", "image/jpg"];
    if (!allowedTypes.includes(file.type)) {
        showAlertBox({
            title: "Invalid File",
            message: "Please select a valid image file (JPG, PNG, WEBP)",
            type: "danger",
            sound: "/assets/sounds/error.mp3"
        });
        return;
    }

    if (file.size > 2 * 1024 * 1024) {
        showAlertBox({
            title: "File Too Large",
            message: "Profile photo size must be less than 2MB",
            type: "danger",
            sound: "/assets/sounds/error.mp3"
        });
        return;
    }

    let confirmed = await confirmDialog("Upload Photo", "Are you sure you want to update your profile photo?");
    if (!confirmed) {
        $(input).val("");
        return;
    }

    let formData = new FormData();
    formData.append("photo", file);

    let uploadBtn = $("#btn-trigger-photo-upload");
    uploadBtn.find("i").removeClass("ti-camera").addClass("ti-loader");

    try {
        let response = await web.post("/profile/change-photo", formData);
        if (response.data && response.data.data && response.data.data.url) {
            let newUrl = response.data.data.url;
            $("#profile-card-avatar").attr("src", newUrl);
            $(".header .avatar img, .sidebar .avatar img, .two-col-sidebar .avatar img").attr("src", newUrl);

            showAlertBox({
                title: "Success",
                message: "Profile Photo Updated Successfully",
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });
        }
    } catch (err) {
        showAlertBox({
            title: "Error",
            message: "Failed to upload profile photo",
            type: "danger",
            sound: "/assets/sounds/error.mp3"
        });
    } finally {
        uploadBtn.find("i").removeClass("ti-loader").addClass("ti-camera");
        $(input).val("");
    }
}

function switchProfileTab(targetId) {

    if (!targetId) return;

    targetId = String(targetId).replace("#", "");

    // 1. Update active tab indicator on top nav tabs
    $("#profile-top-tabs .nav-link").removeClass("active").attr("aria-selected", "false");

    let activeTab = $(`#profile-top-tabs a[data-target-pane="${targetId}"], #profile-top-tabs a[data-bs-target="#${targetId}"], #profile-top-tabs a[href="#${targetId}"]`);
    if (activeTab.length > 0) {
        activeTab.addClass("active").attr("aria-selected", "true");

        if (typeof bootstrap !== 'undefined' && bootstrap.Tab) {
            try {
                let bsTab = bootstrap.Tab.getOrCreateInstance(activeTab[0]);
                bsTab.show();
            } catch (e) {
                // Gracefully ignore if Bootstrap transition is in progress
            }
        }
    }

    // 2. Hide all tab panes
    $("#profile-tab-content .tab-pane").removeClass("show active").hide();

    // 3. Display target tab pane
    let targetPane = $("#" + targetId);
    if (targetPane.length > 0) {
        targetPane.show();
        if (targetPane[0]) {
            void targetPane[0].offsetHeight; // force reflow for smooth display
        }
        targetPane.addClass("show active");
    }
}

window.destroyPage = function () {

    $(document).off(".profileTabSwitch");
    $("#form-edit-profile").off(".profileSettings");
    $("#form-change-password").off(".passwordSettings");
    $("#profile-wrapper").off(".profileContainer");
};
