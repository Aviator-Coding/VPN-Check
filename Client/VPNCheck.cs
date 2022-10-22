using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using Microsoft.VisualBasic;

namespace VPN_Check
{

    public enum Block
    {
        Residential,
        Proxy,
        Warning
    }

    static class BlockExtensions
    {
        public static bool isProxy(this Block block)
        {
            if (block == Block.Proxy)
                return true;
            else return false;
        }
    }


//     {
//     "Message": "success",
//     "Date": "2022-10-22T21:04:07.876673367Z",
//     "Success": true,
//     "Error": "",
//     "Data": {
//         "IP": "",
//         "CountryCode": "",
//         "CountryName": "",
//         "Asn": 0,
//         "Isp": "",
//         "Block": 0
//     }
// }
    public class ApiResponse
    {
        public string? Message { get; set; }
        public DateTime Date { get; set; }
        public bool Success { get; set; }
        public string? Error { get; set; }
        public ApiData? Data { get; set; }
    }
    
    public class ApiData
    {
        public string? IP { get; set; }
        public string? CountryCode { get; set; }
        public string? CountryName { get; set; }
        public int Asn { get; set; }
        public string Isp { get; set; }
        public Block Block { get; set; }
    }



    public class VPNCheck
    {
        private const string APIKEY = "cXJfTVRfbVpxbUJLY1dFLUhqVHVsTGM5UFhkRGlnS21JNUNLX0hHS0dOOER2LUQ4LU9iRXNlX3ZULVpFY0YzYnBBbkMzdzBlcnZ1b1hZMGY=";
        private const string APIREQUEST = "aHR0cDovL3ZwbmNoZWNrLmF2aWF0b3ItY29kaW5nLmRlL2FwaQ==";
        
        public static async Task<ApiResponse?> CheckVpn()
        {
            using var client = new HttpClient();
            try
            {
                var webRequest = new HttpRequestMessage(HttpMethod.Get, Encoding.UTF8.GetString(Convert.FromBase64String(APIREQUEST)));
                webRequest.Headers.Add("X-Key", Encoding.UTF8.GetString(Convert.FromBase64String(APIKEY)));

                using (var r = await client.SendAsync(webRequest))
                {
                    var response = await r.Content.ReadAsStringAsync();
                    // Debug.WriteLine(response);
                    var result = JsonSerializer.Deserialize<ApiResponse>(response);
                    // Debug.WriteLine(JsonSerializer.Serialize(result));
                    return JsonSerializer.Deserialize<ApiResponse>(response);
                }
            }
            catch (Exception e)
            {
                Debug.WriteLine(e.Message);
                return new ApiResponse();
            }
        }

    }
}