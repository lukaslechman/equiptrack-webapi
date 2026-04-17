const mongoHost = process.env.EQUIPTRACK_API_MONGODB_HOST
const mongoPort = process.env.EQUIPTRACK_API_MONGODB_PORT

const mongoUser = process.env.AMBULANCE_API_MONGODB_USERNAME
const mongoPassword = process.env.AMBULANCE_API_MONGODB_PASSWORD

const database = process.env.EQUIPTRACK_API_MONGODB_DATABASE
const collection = process.env.EQUIPTRACK_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
      print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })

//insert sample data
let result = db[collection].insertMany([
    {
        "id": "eq-001",
        "name": "Röntgen EVO-3000",
        "category": "imaging",
        "manufacturer": "Siemens",
        "serialNumber": "RX-2021-00123",
        "purchaseDate": "2021-03-15",
        "warrantyUntil": "2024-03-15",
        "lifespanYears": 10,
        "purchasePrice": 85000,
        "status": "active",
        "note": "",
        "location": {
            "building": "Budova A",
            "department": "Rádiológia",
            "room": "M-12"
        }
    },
    {
        "id": "eq-002",
        "name": "EKG Monitor CardioLife",
        "category": "monitoring",
        "manufacturer": "Philips",
        "serialNumber": "CL-2019-00456",
        "purchaseDate": "2019-07-20",
        "warrantyUntil": "2022-07-20",
        "lifespanYears": 8,
        "purchasePrice": 12000,
        "status": "damaged",
        "note": "Poškodený displej",
        "location": {
            "building": "Budova B",
            "department": "Kardiológia",
            "room": "M-03"
        }
    },
    {
        "id": "eq-003",
        "name": "Infúzna pumpa InfuMed",
        "category": "infusion",
        "manufacturer": "B.Braun",
        "serialNumber": "IM-2020-00789",
        "purchaseDate": "2020-05-10",
        "warrantyUntil": "2023-05-10",
        "lifespanYears": 7,
        "purchasePrice": 4500,
        "status": "active",
        "note": "",
        "location": {
            "building": "Budova A",
            "department": "Chirurgia",
            "room": "M-07"
        }
    },
    {
        "id": "eq-004",
        "name": "Defibrilátor ShockPro",
        "category": "emergency",
        "manufacturer": "Zoll",
        "serialNumber": "SP-2018-00321",
        "purchaseDate": "2018-11-01",
        "warrantyUntil": "2021-11-01",
        "lifespanYears": 6,
        "purchasePrice": 9800,
        "status": "decommissioned",
        "note": "Vyradený po kalibrácii",
        "location": {
            "building": "Budova C",
            "department": "OAIM",
            "room": "M-01"
        }
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);