// @ts-nocheck
export async function NFCWriteURL(url) {
    const ndef = new NDEFReader();
    await ndef.write({
        records: [{ recordType: "url", data: url }]
    });
}

export async function NFCScan(log) {
    const ndef = new NDEFReader();
    ndef.addEventListener("readingerror", () => {
        log("Argh! Cannot read data from the NFC tag. Try another one?");
    });

    ndef.addEventListener("reading", ({ message, serialNumber }) => {
        log(`> Serial Number: ${serialNumber}`);
        log(`> Records: (${message.records.length})`);
        const decoder = new TextDecoder();
        for (const record of message.records) {
            log("Record type:  " + record.recordType);
            log("MIME type:    " + record.mediaType);
            log("=== data ===\n" + decoder.decode(record.data));
        }
    });
    await ndef.scan();
}