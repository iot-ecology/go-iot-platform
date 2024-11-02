import 'ol/ol.css';
import Map from 'ol/Map';
import View from 'ol/View';
import TileLayer from 'ol/layer/Tile';
import OSM from 'ol/source/OSM';
import {fromLonLat} from 'ol/proj';
import {Card} from 'antd';
import React, {useEffect} from "react";
import {Feature} from "ol";
import VectorSource from "ol/source/Vector";
import VectorLayer from "ol/layer/Vector";
import {Point} from "ol/geom";
import {Circle, Fill, Stroke, Style} from "ol/style";


const DeviceMap: React.FC = () => {
  useEffect(() => {
    const map = new Map({
      target: 'map', // 地图容器ID
      layers: [new TileLayer({
        source: new OSM(),
      }),], view: new View({
        center: fromLonLat([116.4074, 39.9042]), // 北京的经纬度
        zoom: 10,
      }), controls: [], // 移除所有控制
    });

// 添加设备标记示例
    const deviceCoordinates = [{lon: 116.4074, lat: 39.9042}, // 示例设备坐标
      {lon: 116.5, lat: 39.9}, {lon: 116.4, lat: 39.92},];

    const features = deviceCoordinates.map(device => {
      var feature = new Feature({
        geometry: new Point(fromLonLat([device.lon, device.lat])),
      });

      // 设置样式
      feature.setStyle(new Style({
        image: new Circle({
          radius: 5, // 半径
          fill: new Fill({color: 'red'}), // 填充颜色
          stroke: new Stroke({color: 'white', width: 2}), // 边框颜色和宽度
        }),
      }));

      return feature;
    });

    const vectorSource = new VectorSource({
      features: features,
    });

    const vectorLayer = new VectorLayer({
      source: vectorSource,
    });

    map.addLayer(vectorLayer); // 将设备标记图层添加到地图

    return () => map.setTarget(undefined); // 清理
  }, []);

  return (<Card title={"Device Map"}>
      <div id="map" style={{width: '100%', height: '690px'}}/>
    </Card>);
};
export default DeviceMap;
