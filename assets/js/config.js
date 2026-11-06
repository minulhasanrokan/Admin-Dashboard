'use strict';

let BaseUrl = window.location.origin;

let activeRequests = 0;
let loader = document.getElementById('loaderOverlay');

function showLoader() {
    if (activeRequests === 0 && loader) {
        loader.classList.remove('hidden');
    }
    activeRequests++;
}

function hideLoader() {
    activeRequests--;
    if (activeRequests <= 0) {
        activeRequests = 0;
        setTimeout(() => {
            if (activeRequests === 0 && loader) {
                loader.classList.add('hidden');
            }
        }, 300);
    }
}

window.addEventListener('load', () => {
    activeRequests = 0;
    if(loader) loader.classList.add('hidden');
});

let web = axios.create({
    baseURL: BaseUrl,
    timeout: 30 * 60 * 1000,
    headers: {
        "Content-Type": "application/json",
        "Accept": "application/json",
    }
});

web.interceptors.request.use(

    function (config) {
        showLoader();

        let _csrf = document.querySelector('meta[name="_csrf"]')?.getAttribute('content');

        if (_csrf) {
            config.headers['X-Csrf-Token'] = _csrf;
        }
        let token = localStorage.getItem("token");

        if (token) {
            config.headers.Authorization = "Bearer " + token;
        }
        if (config.data instanceof FormData) {
            delete config.headers['Content-Type'];
        }

        return config;
    },
    function (error) {
        hideLoader();
        return Promise.reject(error);
    }
);

web.interceptors.response.use(
    function (response) {
        let newCsrf = response.headers['x-csrf-token'] || response.headers['X-Csrf-Token'] || response.data?._csrf || response.data?.CSRF;

        if (newCsrf) {
            document.querySelector('meta[name="_csrf"]')?.setAttribute('content', newCsrf);
        }

        hideLoader();

        return response;
    },
    function (error) {
        let newCsrf = error.response?.headers['x-csrf-token'] || error.response?.headers['X-Csrf-Token'] || error.response?.data?._csrf || error.response?.data?.CSRF;

        if (newCsrf) {
            document.querySelector('meta[name="_csrf"]')?.setAttribute('content', newCsrf);
        }

        hideLoader();

        if (error.response) {
            let status = error.response.status;
            switch (status) {
                case 400:
                    showAlertBox({title: "Bad Request", message: error.response.data.message || "Invalid Request", type: "danger", sound: "/assets/sounds/error.mp3"});
                    break;
                case 401:
                    showAlertBox({title: "Unauthorized", message: "Please Login Again", type: "danger", sound: "/assets/sounds/error.mp3"});
                    break;
                case 403:
                    showAlertBox({title: "Forbidden", message: "You Don’t Have Permission", type: "danger", sound: "/assets/sounds/error.mp3"});
                    break;
                case 404:
                    showAlertBox({title: "Not Found", message: "Resource Not Found", type: "danger", sound: "/assets/sounds/error.mp3"});
                    break;
                case 422:

                    let message = error.response?.data?.message || "Validation Failed";

                    if (typeof parseErrors === 'function' && typeof showErrorMsg === 'function') {
                        let errorData = parseErrors(error.response.data);
                        Object.keys(errorData).forEach(field => {
                            showErrorMsg(field, errorData[field]);
                        });
                    }

                    showAlertBox({title: "Failed", message: message, type: "danger", sound: "/assets/sounds/error.mp3"});
                    break;
                case 500:
                    showAlertBox({title: "Server Error", message: "Something went wrong", type: "danger", sound: "/assets/sounds/error.mp3"});
                    break;
                default:
                    showAlertBox({title: "Error", message: error.response.data.message || "Something Went Wrong", type: "danger", sound: "/assets/sounds/error.mp3"});
            }
        }
        else if (error.request) {
            showAlertBox({title: "Network Error", message: "No Response From Server", type: "danger", sound: "/assets/sounds/error.mp3"});
        }
        else {

            showAlertBox({title: "Request Error", message: "Request Error", type: "danger", sound: "/assets/sounds/error.mp3"});
        }

        hideLoader();

        return Promise.reject(error);
    }
);

function showHideValue(id, button_id) {

    let input = document.getElementById(id);
    let btn = document.getElementById(button_id);

    if (input.type === "password") {
        input.type = "text";
        btn.className = "ti ti-eye";
    }
    else {
        input.type = "password";
        btn.className = "ti ti-eye-off";
    }
}
let containerAssets = new Map();


function loadCSS(href) {

    if (!href) {
        return Promise.resolve();
    }

    return new Promise((resolve) => {

        let existing = document.querySelector(
            `link[data-dyn-asset="${CSS.escape(href)}"]`
        );

        if (existing) {
            resolve();
            return;
        }

        let link = document.createElement('link');

        link.rel = 'stylesheet';
        link.href = href;
        link.dataset.dynAsset = href;

        link.onload = function () {
            resolve();
        };

        link.onerror = function () {
            resolve();
        };

        document.head.appendChild(link);
    });
}

function loadScript(src) {

    if (!src) {
        return Promise.resolve();
    }

    return new Promise((resolve, reject) => {

        let existing = document.querySelector(
            `script[data-dyn-asset="${CSS.escape(src)}"]`
        );

        if (existing) {
            resolve();
            return;
        }

        let script = document.createElement('script');

        script.src = src;
        script.dataset.dynAsset = src;

        script.onload = function () {
            resolve();
        };

        script.onerror = function () {

            reject(
                new Error(
                    'Failed to load script: ' + src
                )
            );
        };

        document.head.appendChild(script);
    });
}

function unloadContainerAssets(containerKey) {

    let prev = containerAssets.get(containerKey);

    if (!prev) {
        return;
    }

    if (typeof window.destroyPage === 'function') {

        window.destroyPage();
        window.destroyPage = null;
    }

    if (Array.isArray(prev.css)) {

        prev.css.forEach(function (href) {

            if (typeof href !== 'string') {
                return;
            }

            let selector = `link[data-dyn-asset="${CSS.escape(href)}"]`;

            document.querySelectorAll(selector).forEach(function (el) {el.onload = null;el.onerror = null;el.remove();});
        });
    }

    if (Array.isArray(prev.js)) {

        prev.js.forEach(function (src) {

            if (typeof src !== 'string') {
                return;
            }

            let selector = `script[data-dyn-asset="${CSS.escape(src)}"]`;

            document.querySelectorAll(selector).forEach(function (el) {el.onload = null;el.onerror = null;el.remove();});
        });
    }

    containerAssets.delete(containerKey);
}

function extractAndLoadAssets(html) {

    let temp = document.createElement('div');

    temp.innerHTML = html;

    let cssUrls = [...temp.querySelectorAll('link[rel="stylesheet"]')].map(function (el) {return el.getAttribute('href');}).filter(Boolean);

    let jsUrls = [...temp.querySelectorAll('script[src]')].map(function (el) {return el.getAttribute('src');}).filter(Boolean);

    temp.querySelectorAll('link[rel="stylesheet"], script[src]').forEach(function (el) {el.remove();});

    return {cleanedHtml: temp.innerHTML, cssUrls: cssUrls, jsUrls: jsUrls};
}

async function loadContainerAssets(containerKey, cssUrls = [], jsUrls = []) {

    if (containerAssets.has(containerKey)) {
        unloadContainerAssets(containerKey);
    }

    let assets = {css: [], js: []};

    for (let href of cssUrls) {

        await loadCSS(href);
        assets.css.push(href);
    }

    for (let src of jsUrls) {

        await loadScript(src);
        assets.js.push(src);
    }
}