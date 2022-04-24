using System;

using System.Runtime.InteropServices;

using System.IO.Compression;

using System.IO;

using System.Diagnostics;

using System.Security.Principal;

using System.Xml.Linq;

using System.Linq;

//using Microsoft.AspNetCore.Builder;

//using System.Net.WebSockets;

//using System.Threading;

using Newtonsoft.Json;

using WebSocketSharp;


namespace CoopirAgent
{
    class Program
    {

        static void Main(string[] args)
        {

            //String s = String.Format("We are running on {0}", RuntimeInformation.OSDescription.ToString());

            //Console.WriteLine(s);

            Guid myuuid = Guid.NewGuid();
            string myuuidAsString = myuuid.ToString();

            MachineInfo machine = new MachineInfo()
            {
                uuid = myuuidAsString,
                name = Environment.MachineName,
                os = Environment.OSVersion.ToString(),
                arch = RuntimeInformation.ProcessArchitecture.ToString()
            };

            using (WebSocket ws = new WebSocket("ws://localhost:4201"))
            {
                ws.OnMessage += Ws_OnMessage;

                ws.Connect();
                string stringjson = JsonConvert.SerializeObject(machine);
                Console.WriteLine(stringjson);
                ws.Send(stringjson);

                Console.ReadKey();

                ws.Close();
            }

            //Zipper();

        }

        private static void Ws_OnMessage(object sender, MessageEventArgs e)
        {
            Console.WriteLine("Recieved from the server: " + e.Data.ToString());

            //Zipper();
            //byte[] zipfile;
            //if (OperatingSystem.IsWindows())
            //    zipfile = File.ReadAllBytes(@".\zip\*.zip");
            //else
            //    zipfile = File.ReadAllBytes("@./zip/*.zip");
            //WebSocket ws = (WebSocket)sender;
            //ws.SendAsync(zipfile);
        }

        static void Zipper()
        {

            if (OperatingSystem.IsWindows() == true)
            {
                String s = String.Format("{0} detected attempting log extraction...", RuntimeInformation.OSDescription.ToString());
                Console.WriteLine(s);
                try
                {
                    if (!IsAdministrator())
                    {
                        ExecuteAsAdmin();
                        Environment.Exit(0);
                    }

                    EventLog[] eventLogs;

                    eventLogs = EventLog.GetEventLogs(Environment.MachineName);

                    String message = String.Format("Number of logs on {0}: ", Environment.MachineName);
                    Console.WriteLine(message + eventLogs.Length);

                    if (!Directory.Exists(@".\Logs"))
                    {
                        DirectoryInfo di = Directory.CreateDirectory(@".\Logs");
                    }

                    foreach (EventLog log in eventLogs)
                    {
                        Console.WriteLine("Log: " + log.Log);

                        var xml = new XDocument(
                        new XElement(log.LogDisplayName.Replace(' ', '_'),
                            from EventLogEntry entry in log.Entries
                            orderby entry.TimeGenerated descending
                            select new XElement("Log",
                              new XElement("Message", entry.Message),
                              new XElement("TimeGenerated", entry.TimeGenerated),
                              new XElement("Source", entry.Source),
                              new XElement("EntryType", entry.EntryType.ToString())
                            )
                        ));
                        if (!Directory.Exists(@".\Logs\" + log.LogDisplayName + ".xml"))
                            xml.Save(@".\Logs\" + log.LogDisplayName + ".xml");
                        else
                        {
                            File.Delete(@".\Logs\" + log.LogDisplayName + ".xml");
                            xml.Save(@".\Logs\" + log.LogDisplayName + ".xml");
                        }
                    }

                    WindowsZip();
                    Console.WriteLine("Log Extraction Successfully. Press any key to exit...");
                    Console.ReadLine();
                    Environment.Exit(0);
                }
                catch (Exception e)
                {
                    Console.WriteLine(e);
                    Console.WriteLine("Log Extraction Failed. Press any key to exit...");
                    Console.ReadLine();
                    Environment.Exit(0);
                }
            }
            else if (OperatingSystem.IsLinux() == true)
            {
                //Console.WriteLine("Linux OS detected attempting log extraction...");
                String s = String.Format("{0} detected attempting log extraction...", RuntimeInformation.OSDescription.ToString());
                Console.WriteLine(s);
                try
                {
                    try
                    {
                        LinuxZip();
                    }
                    catch (Exception e)
                    {
                        Console.WriteLine("[ERROR]" + e);
                    }
                }
                catch (Exception e)
                {
                    Console.WriteLine(e);
                }
            }
            else if (OperatingSystem.IsMacOS() == true)
            {
                String s = String.Format("{0} detected attempting log extraction...", RuntimeInformation.OSDescription.ToString());
                Console.WriteLine(s);
                try
                {
                    try
                    {
                        MacZip();
                    }
                    catch (Exception e)
                    {
                        Console.WriteLine("[ERROR]" + e);
                    }
                }
                catch (Exception e)
                {
                    Console.WriteLine(e);
                }
            }
            else
            {
                Console.WriteLine("[Error] Cannot determine OS");
            }
        }

        public static bool IsAdministrator()
        {
            var identity = WindowsIdentity.GetCurrent();
            var principal = new WindowsPrincipal(identity);
            return principal.IsInRole(WindowsBuiltInRole.Administrator);
        }

        static void ExecuteAsAdmin()
        {
            Process proc = new Process();
            proc.StartInfo.FileName = "CoopirAgent.exe";
            proc.StartInfo.UseShellExecute = true;
            proc.StartInfo.Verb = "runas";
            proc.Start();
        }

        static void WindowsZip()
        {
            string startFolder = @".\Logs";
            string zipFile = string.Format(@".\zip\evtLogs_{0}.zip", Environment.MachineName);
            string zipFolder = @".\zip";

            if (!Directory.Exists(zipFolder))
            {
                Console.WriteLine(@"Making directory .\Zip");
                DirectoryInfo di = Directory.CreateDirectory(zipFolder);
            }
            
            try
            {
                ZipFile.CreateFromDirectory(startFolder, zipFile);
            }
            catch(System.IO.IOException e)
            {
                Console.WriteLine("Removing Older Zip");
                File.Delete(zipFile);
                ZipFile.CreateFromDirectory(startFolder, zipFile);
            }
        }

        static void LinuxZip()
        {
            string startFolder = @"/var/log";
            string zipFile = string.Format(@"./zip/Logs_{0}.zip", Environment.MachineName);
            string zipFolder = @"./zip";
            //string extractFolder = @".\extracted";

            if (!Directory.Exists(zipFolder))
            {
                Console.WriteLine(@"Making directory .\Zip");
                DirectoryInfo di = Directory.CreateDirectory(zipFolder);
                Console.WriteLine("Directory creation successful");
            }

            try
            {
                String s = String.Format("Attempting Zip of {0}", startFolder);
                Console.WriteLine(s);
                ZipFile.CreateFromDirectory(startFolder, zipFile);
                Console.WriteLine("Zip Successful");

            }
            catch (System.IO.IOException e)
            {
                Console.WriteLine("Removing Older Zip");
                File.Delete(zipFile);
                Console.WriteLine("Zip Removed");
                String s = String.Format("Attempting Zip of {0}", startFolder);
                Console.WriteLine(s);
                ZipFile.CreateFromDirectory(startFolder, zipFile);
                Console.WriteLine("Zip Successful");
            }
        }

        static void MacZip()
        {
            string startFolder = @"/var/log";
            string zipFile = string.Format(@"./zip/Logs_{0}.zip", Environment.MachineName);
            string zipFolder = @".\zip";
            //string extractFolder = @".\extracted";

            if (!Directory.Exists(zipFolder))
            {
                Console.WriteLine(@"Making directory .\Zip");
                DirectoryInfo di = Directory.CreateDirectory(zipFolder);
            }

            try
            {
                ZipFile.CreateFromDirectory(startFolder, zipFile);
            }
            catch (System.IO.IOException e)
            {
                Console.WriteLine("Removing Older Zip");
                File.Delete(zipFile);
                ZipFile.CreateFromDirectory(startFolder, zipFile);
            }
        }


    }

    public static class OperatingSystem
    {
        public static bool IsWindows() =>
            RuntimeInformation.IsOSPlatform(OSPlatform.Windows);

        public static bool IsMacOS() =>
            RuntimeInformation.IsOSPlatform(OSPlatform.OSX);

        public static bool IsLinux() =>
            RuntimeInformation.IsOSPlatform(OSPlatform.Linux);


    }

    public class MachineInfo
    {
        public string uuid { get; set; }
        public string name { get; set; }
        public string os { get; set; }
        public string arch { get; set; }
    }


}
