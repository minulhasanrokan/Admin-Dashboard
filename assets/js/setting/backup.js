'use strict';

(function ($) {
    let backupGrid = null;

    $(document).ready(async function () {

        let backup = $("#backup_type").val();

        await loadDropDown('backup_types', 'code,name', 'backup', 'backup_container', 'Select Backup', 0, backup, 0, "8", "0", "0", 0, 0);

        bindBackupSettingsEvents();
        initBackupHistoryGrid();
    });

    function bindBackupSettingsEvents() {

        let form = $("#backup-settings-form");

        form.off(".backupSettings");

        form.on("focusout.backupSettings", "select", function () { validateField(this.id); });

        form.on("select2:close.backupSettings", "select", function () { validateField(this.id); });

        form.on("click.backupSettings", "#backup_setting_button", async function () {
            await saveBackupSettings();
        });

        $(document).off("click.backupDownload", ".btn-backup-download").on("click.backupDownload", ".btn-backup-download", function (e) {
            e.preventDefault();
            let uuid = $(this).attr("data-uuid") || $(this).attr("data-id");
            if (!uuid) return;

            let downloadUrl = `/settings/backup/download/${uuid}`;
            let a = document.createElement('a');
            a.href = downloadUrl;
            a.download = '';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);

            setTimeout(function () {
                if (backupGrid) {
                    backupGrid.reload();
                }
            }, 800);
        });

        $(document).off("click.backupError", ".view-backup-error").on("click.backupError", ".view-backup-error", function (e) {
            e.preventDefault();
            let errMsg = $(this).attr("data-error") || "No error details available";
            showAlertBox({
                title: "Backup Error Details",
                message: errMsg,
                type: "danger",
                sound: "/assets/sounds/error.mp3"
            });
        });
    }

    function initBackupHistoryGrid() {
        backupGrid = initAgDataGrid({
            gridContainer: '#backup_ag_grid',
            url: '/settings/backup/list',
            exportFileName: 'backup_history',
            defaultOrderCol: 'created_at',
            defaultOrderDir: 'desc',
            totalBadgeSelector: '#backup_total_badge',
            columnDefs: [
                {
                    headerName: '#',
                    isIndexCol: true,
                    width: 70,
                    minWidth: 65
                },
                {
                    field: 'backup_type',
                    colId: 'backup_type',
                    headerName: 'Backup Type',
                    width: 170,
                    minWidth: 150,
                    sortable: true,
                    filter: true,
                    filterType: 'select',
                    filterOptions: [
                        { value: '0001', label: 'Database' },
                        { value: '0002', label: 'Files' }
                    ],
                    cellRenderer: function (params) {
                        if (params.value === '0002' || (params.data && params.data.backup_type_name === 'Files')) {
                            return '<span class="badge bg-info-transparent text-info fw-semibold"><i class="ti ti-folder me-1"></i>Files</span>';
                        }
                        return '<span class="badge bg-primary-transparent text-primary fw-semibold"><i class="ti ti-database me-1"></i>Database</span>';
                    }
                },
                {
                    field: 'file_size_formatted',
                    colId: 'file_size',
                    headerName: 'File Size',
                    width: 140,
                    minWidth: 130,
                    sortable: true,
                    filter: false,
                    cellRenderer: function (params) {
                        return `<span class="badge bg-light text-secondary border">${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'backup_status',
                    colId: 'backup_status',
                    headerName: 'Status',
                    width: 150,
                    minWidth: 140,
                    sortable: true,
                    filter: true,
                    filterType: 'select',
                    filterOptions: [
                        { value: '1', label: 'Pending' },
                        { value: '2', label: 'Running' },
                        { value: '3', label: 'Completed' },
                        { value: '4', label: 'Failed' }
                    ],
                    cellRenderer: function (params) {
                        let status = params.value;
                        switch (status) {
                            case 3:
                                return '<span class="badge bg-success-transparent text-success fw-semibold"><i class="ti ti-check me-1"></i>Completed</span>';
                            case 2:
                                return '<span class="badge bg-info-transparent text-info fw-semibold"><i class="ti ti-loader-2 me-1"></i>Running</span>';
                            case 4:
                                return '<span class="badge bg-danger-transparent text-danger fw-semibold"><i class="ti ti-alert-triangle me-1"></i>Failed</span>';
                            default:
                                return '<span class="badge bg-warning-transparent text-warning fw-semibold"><i class="ti ti-clock me-1"></i>Pending</span>';
                        }
                    }
                },
                {
                    field: 'backup_time',
                    colId: 'created_at',
                    headerName: 'Backup Time',
                    width: 195,
                    minWidth: 175,
                    sortable: true,
                    filter: true,
                    filterType: 'date',
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-clock me-1"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'download_status',
                    colId: 'download_status',
                    headerName: 'Download Status',
                    width: 160,
                    minWidth: 145,
                    sortable: true,
                    filter: true,
                    filterType: 'select',
                    filterOptions: [
                        { value: '0', label: 'Not Downloaded' },
                        { value: '1', label: 'Downloaded' }
                    ],
                    cellRenderer: function (params) {
                        let status = params.value;
                        if (status === 1 || (params.data && params.data.download_status_name === 'Downloaded')) {
                            return '<span class="badge bg-success-transparent text-success fw-semibold"><i class="ti ti-download me-1"></i>Downloaded</span>';
                        }
                        return '<span class="badge bg-light text-muted border fw-normal"><i class="ti ti-minus me-1"></i>Not Downloaded</span>';
                    }
                },
                {
                    field: 'downloaded_by',
                    colId: 'downloaded_by',
                    headerName: 'Download By',
                    width: 175,
                    minWidth: 150,
                    sortable: true,
                    filter: true,
                    filterType: 'string',
                    cellRenderer: function (params) {
                        let val = params.value || '-';
                        return `<span class="fw-medium text-dark"><i class="ti ti-user me-1 text-muted"></i>${escapeHtml(val)}</span>`;
                    }
                },
                {
                    field: 'download_date',
                    colId: 'download_date',
                    headerName: 'Download Date',
                    width: 195,
                    minWidth: 175,
                    sortable: true,
                    filter: true,
                    filterType: 'date',
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-calendar me-1"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    headerName: 'Action',
                    colId: 'action',
                    width: 110,
                    minWidth: 100,
                    sortable: false,
                    filter: false,
                    cellClass: 'd-flex align-items-center justify-content-center',
                    cellRenderer: function (params) {
                        let row = params.data;
                        if (!row) return '';
                        if (row.backup_status === 3) {
                            if (row.download_status === 1 || row.download_status_name === 'Downloaded') {
                                return '<span class="text-muted fs-12">-</span>';
                            }
                            return `<button type="button" class="btn btn-sm btn-icon btn-outline-primary btn-backup-download" data-uuid="${row.uuid}" data-id="${row.uuid}" title="Download Backup"><i class="ti ti-download"></i></button>`;
                        } else if (row.backup_status === 4) {
                            return `<button type="button" class="btn btn-sm btn-icon btn-outline-danger view-backup-error" data-error="${escapeHtml(row.error_log)}" title="View Error Details"><i class="ti ti-info-circle"></i></button>`;
                        } else {
                            return `<button class="btn btn-sm btn-icon btn-light disabled" disabled><i class="ti ti-clock"></i></button>`;
                        }
                    }
                }
            ]
        });
    }

    function validateField(id) {

        switch (id) {
            case "backup":
                form_validation({control: "backup#select#mandatory", msg_text: "Backup Type"});
                break;
        }
    }

    async function saveBackupSettings() {

        let status = form_validation({control: "backup#select#mandatory", msg_text: "Backup Type"});

        if (status === false) {

            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        let confirmed = await confirmDialog("Generate Backup", "Are you sure you want to generate database backup?");

        if (!confirmed) {
            return;
        }

        let formData = new FormData();

        formData.append("backup", $("#backup").val());

        let response = await web.post('/settings/backup', formData);

        if (response.status === 200) {

            showAlertBox({title: "Success", message: "Backup request generated successfully", type: "info", sound: "/assets/sounds/success.mp3"});
            
            setTimeout(function () {
                if (backupGrid) {
                    backupGrid.reload(1);
                }
            }, 1000);
        }
    }

    function formatFileSize(bytes) {
        if (!bytes || bytes === 0) return "-";
        bytes = Number(bytes);
        if (isNaN(bytes) || bytes <= 0) return "-";

        const KB = 1024;
        const MB = 1024 * KB;
        const GB = 1024 * MB;
        const TB = 1024 * GB;

        if (bytes >= TB) {
            return (bytes / TB).toFixed(2) + " TB";
        } else if (bytes >= GB) {
            return (bytes / GB).toFixed(2) + " GB";
        } else if (bytes >= MB) {
            return (bytes / MB).toFixed(2) + " MB";
        } else if (bytes >= KB) {
            return (bytes / KB).toFixed(2) + " KB";
        } else {
            return bytes + " B";
        }
    }

    function escapeHtml(text) {
        if (!text) return "";
        return String(text)
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    window.destroyPage = function () {

        if (backupGrid) {
            backupGrid.destroy();
            backupGrid = null;
        }

        $("#backup").each(
            function () {

                if (
                    $(this).hasClass(
                        "select2-hidden-accessible"
                    )
                ) {

                    $(this).select2(
                        "destroy"
                    );
                }
            }
        );

        $("#backup-settings-form").off(
            ".backupSettings"
        );

        $(document).off(".backupDownload");
        $(document).off(".backupError");
    };
})(jQuery);