db = db.getSiblingDB("mongodb");

const documents = [
    {
        name: "MICS-6814",
        latitude: -23.562387,
        longitude: -46.711777,
        params: {
            co2: { min: 0, max: 1000, z: 1.96 },
            co: { min: 0, max: 15, z: 1.96 },
            no2: { min: 0, max: 1130, z: 1.96 },
            mp10: { min: 0, max: 250, z: 1.96 },
            mp25: { min: 0, max: 125, z: 1.96 },
            rad: { min: 1, max: 1280, z: 1.96 },
        },
    },
    {
        name: "SPS30",
        latitude: -23.564137,
        longitude: -46.711639,
        params: {
            co2: { min: 0, max: 1000, z: 1.96 },
            co: { min: 0, max: 15, z: 1.96 },
            no2: { min: 0, max: 1130, z: 1.96 },
            mp10: { min: 0, max: 250, z: 1.96 },
            mp25: { min: 0, max: 125, z: 1.96 },
            rad: { min: 1, max: 1280, z: 1.96 },
        },
    },
    {
        name: "RXW-LIB-900",
        latitude: -23.565203,
        longitude: -46.709176,
        params: {
            co2: { min: 0, max: 1000, z: 1.96 },
            co: { min: 0, max: 15, z: 1.96 },
            no2: { min: 0, max: 1130, z: 1.96 },
            mp10: { min: 0, max: 250, z: 1.96 },
            mp25: { min: 0, max: 125, z: 1.96 },
            rad: { min: 1, max: 1280, z: 1.96 },
        },
    },
];

db.sensors.insertMany(documents);

print("Data migration completed successfully. Documents inserted:");
printjson(documents);
