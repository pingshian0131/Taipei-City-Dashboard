<script setup>
import { ref } from "vue";
import DashboardComponent from "../dashboardComponent/DashboardComponent.vue";
import MapContainer from "../components/map/MapContainer.vue";
import { useMapStore } from "../store/mapStore";
import { useDialogStore } from "../store/dialogStore";

const mapStore = useMapStore();
const dialogStore = useDialogStore();

const rainfallComponent = ref({
	id: "flood-rainfall-demo",
	index: "flood_rainfall_demo",
	city: "taipei",
	name: "即時雨量警戒",
	source: "示範資料｜水利處雨量站即時資料 (API #4)",
	time_from: "demo",
	time_to: null,
	update_freq: 10,
	update_freq_unit: "minute",
	chart_config: {
		types: ["ColumnChart"],
		color: ["#5a9cf8", "#f6c344", "#ed5a5a"],
		unit: "mm/hr",
	},
	chart_data: [
		{
			name: "近 1 小時雨量",
			data: [
				{ x: "雙溪站", y: 38 },
				{ x: "天母站", y: 22 },
				{ x: "內湖站", y: 14 },
				{ x: "信義站", y: 9 },
				{ x: "公館站", y: 6 },
			],
		},
	],
	map_config: [null],
});

const pumpStationComponent = ref({
	id: "flood-pump-status-demo",
	index: "flood_pump_status_demo",
	city: "taipei",
	name: "抽水站運轉負載",
	source: "示範資料｜水利處抽水站運轉狀態 (API #5)",
	time_from: "demo",
	time_to: null,
	update_freq: 5,
	update_freq_unit: "minute",
	chart_config: {
		types: ["BarChart"],
		color: ["#5a9cf8"],
		unit: "%",
	},
	chart_data: [
		{
			name: "啟動水泵比例",
			data: [
				{ x: "玉成", y: 80 },
				{ x: "南港", y: 60 },
				{ x: "康樂", y: 40 },
				{ x: "民生", y: 20 },
				{ x: "六館", y: 0 },
			],
		},
	],
	map_config: [
		{
			index: "patrol_rain_floodgate",
			title: "抽水站位置",
			type: "circle",
			source: "geojson",
			size: "small",
			icon: null,
			paint: {
				"circle-color": "#5a9cf8",
				"circle-stroke-color": "#ffffff",
				"circle-stroke-width": 1,
			},
			property: [
				{ key: "station_name", name: "抽水站名稱" },
				{ key: "river_basin", name: "所屬流域" },
				{ key: "all_pumb_lights", name: "目前狀態" },
				{ key: "pumb_num", name: "水泵數量" },
				{ key: "warning_level", name: "警戒水位 (m)" },
			],
			city: "taipei",
		},
	],
});

const floodRadarMapComponent = ref({
	id: "flood-move-car-radar-demo",
	index: "flood_move_car_radar_demo",
	city: "taipei",
	name: "移車雷達：易淹水區 × 開放停車學校",
	source: "示範資料｜易淹水危害圖資 + 颱風期間開放停車學校",
	time_from: "demo",
	time_to: null,
	update_freq: null,
	update_freq_unit: null,
	chart_config: {
		types: ["MapLegend"],
		color: ["#ed5a5a", "#5fcf80"],
		unit: "",
	},
	chart_data: [
		{
			name: "易淹水危害區",
			type: "fill",
			icon: null,
		},
		{
			name: "颱風期間開放停車學校",
			type: "circle",
			icon: null,
		},
	],
	map_config: [
		{
			index: "wee_hazard_water_tp",
			title: "易淹水危害區",
			type: "fill",
			source: "geojson",
			size: null,
			icon: null,
			paint: {
				"fill-color": "#ed5a5a",
				"fill-opacity": 0.35,
			},
			property: [
				{ key: "type", name: "淹水深度 (m)" },
				{ key: "hazard_class", name: "情境分類" },
				{ key: "COUNTY", name: "縣市" },
			],
			city: "taipei",
		},
		{
			index: "flood_typhoon_school_demo",
			title: "颱風期間開放停車學校",
			type: "circle",
			source: "geojson",
			size: "big",
			icon: null,
			paint: {
				"circle-color": "#5fcf80",
				"circle-stroke-color": "#ffffff",
				"circle-stroke-width": 1.5,
			},
			property: [
				{ key: "school_name", name: "學校名稱" },
				{ key: "district", name: "行政區" },
				{ key: "parking_capacity", name: "可停車位" },
				{ key: "contact", name: "聯絡電話" },
			],
			city: "taipei",
		},
	],
});

const toggleOn = ref({
	rainfall: false,
	pump: false,
	radar: false,
});

function shouldDisable(map_config) {
	if (!map_config?.[0]) return true;
	const allMapLayerIds = map_config.map(
		(el) => `${el.index}-${el.type}-${el.city}`,
	);
	if (mapStore.isPreloading === true) {
		return true;
	}
	return (
		mapStore.loadingLayers.filter((el) => allMapLayerIds.includes(el))
			.length > 0
	);
}

function handleToggle(value, map_config, key) {
	toggleOn.value[key] = value;
	if (!map_config?.[0]) {
		if (value) {
			dialogStore.showNotification(
				"info",
				"本組件沒有空間資料，不會渲染地圖",
			);
		}
		return;
	}
	if (value) {
		mapStore.addToMapLayerList(map_config);
	} else {
		mapStore.clearByParamFilter(map_config);
		mapStore.turnOffMapLayerVisibility(map_config);
	}
}

</script>

<template>
  <div class="floodradarview">
    <header class="floodradarview-header">
      <div>
        <h1>城市淹水預言家 × 移車雷達</h1>
        <p>韌性防災 demo｜整合即時雨量、抽水站運轉、易淹水區、開放停車學校</p>
      </div>
      <RouterLink
        class="floodradarview-header-back"
        to="/dashboard"
      >
        <span>arrow_back</span>
        返回主儀表板
      </RouterLink>
    </header>
    <div class="floodradarview-body">
      <div class="floodradarview-body-charts">
        <h2>淹水預言家</h2>
        <DashboardComponent
          :config="rainfallComponent"
          mode="map"
          :info-btn="false"
          :toggle-disable="shouldDisable(rainfallComponent.map_config)"
          :toggle-on="toggleOn.rainfall"
          @toggle="
            (value, map_config) => handleToggle(value, map_config, 'rainfall')
          "
        />
        <DashboardComponent
          :config="pumpStationComponent"
          mode="map"
          :info-btn="false"
          :toggle-disable="shouldDisable(pumpStationComponent.map_config)"
          :toggle-on="toggleOn.pump"
          @toggle="
            (value, map_config) => handleToggle(value, map_config, 'pump')
          "
        />
        <h2>移車雷達</h2>
        <DashboardComponent
          :config="floodRadarMapComponent"
          mode="map"
          :info-btn="false"
          :toggle-disable="shouldDisable(floodRadarMapComponent.map_config)"
          :toggle-on="toggleOn.radar"
          @toggle="
            (value, map_config) => handleToggle(value, map_config, 'radar')
          "
        />
        <p class="floodradarview-body-charts-tip">
          打開組件右上角的開關，可在右側地圖疊加對應圖層。
        </p>
      </div>
      <MapContainer />
    </div>
  </div>
</template>

<style scoped lang="scss">
.floodradarview {
	width: 100vw;
	height: 100vh;
	height: calc(var(--vh) * 100);
	display: flex;
	flex-direction: column;
	background-color: var(--color-background);

	&-header {
		height: 60px;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0 var(--font-l);
		border-bottom: solid 1px var(--color-border);
		background-color: var(--color-component-background);

		h1 {
			color: var(--color-normal-text);
			font-size: var(--font-l);
			font-weight: 500;
		}

		p {
			color: var(--color-complement-text);
			font-size: var(--font-s);
		}

		&-back {
			display: flex;
			align-items: center;
			gap: 4px;
			padding: 6px 12px;
			border-radius: 5px;
			background-color: var(--color-background);
			color: var(--color-highlight);
			font-size: var(--font-s);
			text-decoration: none;
			transition: opacity 0.2s;

			&:hover {
				opacity: 0.8;
			}

			span {
				font-family: var(--font-icon);
				font-size: 1rem;
			}
		}
	}

	&-body {
		flex: 1;
		display: flex;
		margin: var(--font-m);
		overflow: hidden;

		&-charts {
			width: 360px;
			max-height: 100%;
			display: flex;
			flex-direction: column;
			row-gap: var(--font-m);
			margin-right: var(--font-s);
			border-radius: 5px;
			overflow-y: scroll;

			@media (min-width: 1000px) {
				width: 370px;
			}

			@media (min-width: 2000px) {
				width: 400px;
			}

			h2 {
				margin: 4px 0 0 4px;
				color: var(--color-complement-text);
				font-size: var(--font-m);
				font-weight: 500;
			}

			&-tip {
				padding: var(--font-s);
				color: var(--color-complement-text);
				font-size: var(--font-s);
			}
		}
	}
}
</style>
