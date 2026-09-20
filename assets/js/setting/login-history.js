'use strict';

(function ($) {
    let loginGrid = null;

    $(document).ready(function () {
        initLoginHistoryGrid();
    });

    function initLoginHistoryGrid() {
        loginGrid = initAgDataGrid({
            gridContainer: '#login_ag_grid',
            url: '/settings/login-history/list',
            exportFileName: 'user_login_history',
            defaultOrderCol: 'created_at',
            defaultOrderDir: 'desc',
            totalBadgeSelector: '#login_total_badge',
            columnDefs: [
                {
                    headerName: '#',
                    isIndexCol: true,
                    width: 70,
                    minWidth: 65
                },
                {
                    field: 'user_id',
                    colId: 'user_id',
                    headerName: 'User ID',
                    width: 120,
                    minWidth: 100,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="badge bg-primary-transparent text-primary fw-bold">#${params.value || '0'}</span>`;
                    }
                },
                {
                    field: 'ip_address',
                    colId: 'ip_address',
                    headerName: 'Client IP',
                    width: 160,
                    minWidth: 140,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="text-dark fw-medium fs-13"><i class="ti ti-map-pin me-1 text-primary"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'user_agent',
                    colId: 'user_agent',
                    headerName: 'Browser & Platform',
                    width: 280,
                    minWidth: 220,
                    sortable: false,
                    cellRenderer: function (params) {
                        return `<span class="text-truncate d-inline-block mw-100 fs-12 text-muted" title="${params.value || ''}">${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'status_text',
                    colId: 'log_out_status',
                    headerName: 'Session State',
                    width: 150,
                    minWidth: 130,
                    sortable: true,
                    filter: true,
                    filterType: 'select',
                    filterOptions: [
                        { value: '0', label: 'Active' },
                        { value: '1', label: 'Logged Out' }
                    ],
                    cellRenderer: function (params) {
                        let isOut = params.data && params.data.log_out_status !== 0;
                        if (isOut) {
                            return '<span class="badge bg-secondary-transparent text-secondary fw-semibold"><i class="ti ti-power me-1"></i>Logged Out</span>';
                        }
                        return '<span class="badge bg-success-transparent text-success fw-semibold"><i class="ti ti-circle-filled fs-8 me-1"></i>Active</span>';
                    }
                },
                {
                    field: 'log_in_time',
                    colId: 'log_in_time',
                    headerName: 'Login Timestamp',
                    width: 200,
                    minWidth: 180,
                    sortable: true,
                    filter: true,
                    filterType: 'date',
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-login me-1 text-success"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'created_at',
                    colId: 'created_at',
                    headerName: 'Created At',
                    width: 200,
                    minWidth: 180,
                    sortable: true,
                    filter: true,
                    filterType: 'date',
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-clock me-1"></i>${params.value || '-'}</span>`;
                    }
                }
            ]
        });
    }
})(jQuery);
