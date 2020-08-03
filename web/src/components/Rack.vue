<template>
  <div class="col-12 col-xl-6 col-xxl-4 rack">
    <div class="card mb-4 border-bottom-0">
      <div class="card-header">
        <h4 class="mb-0 font-weight-bold text-secondary d-flex">
          <div class="flex-grow-1">{{ rack.name }}</div>
          <div class="flex-shrink-0">
            <i class="fa fa-check-square text-success"></i>
          </div>
        </h4>
      </div>
      <ul class="list-group list-group-flush">
        <li class="list-group-item d-flex align-items-center p-0">
          <div class="flex-even p-3 border-right">
            <div class="font-weight-bold">Apps</div>
            <div>{{ rack.apps.length }}</div>
          </div>
          <div class="flex-even p-3 border-right">
            <div class="font-weight-bold">CPU</div>
            <div>3784 / 4096</div>
          </div>
          <div class="flex-even p-3">
            <div class="font-weight-bold">Memory</div>
            <div>12.4 / 16.0 GB</div>
          </div>
        </li>
        <li class="list-group-item p-0">
          <div class="row" style="padding-left: 15px; padding-right: 14px;">
            <div class="col-12 col-xxl-6 p-3 border-right border-bottom bg-light">
              <div
                style="width: 100%; background-color: #fff; height: 80px; border: 1px #eee solid;"
                class="d-flex align-items-center justify-content-center text-secondary"
              >CPU/Memory Graph</div>
            </div>
            <div class="col-12 col-xxl-6 p-3 border-right border-bottom bg-light">
              <div
                style="width: 100%; background-color: #fff; height: 80px; border: 1px #eee solid;"
                class="d-flex align-items-center justify-content-center text-secondary"
              >Network Graph</div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: ["rack"]
};
</script>
