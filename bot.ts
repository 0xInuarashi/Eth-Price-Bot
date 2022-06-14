import fetch from "node-fetch";
import Discord from "discord.js";
import { ActivityTypes } from "discord.js/typings/enums";

const client = new Discord.Client({
  shards: "auto",
  intents: [Discord.Intents.FLAGS.GUILDS],
});

client.login(process.env.TOKEN);

async function getPrice() {
  const raw = await fetch("https://api.coinbase.com/v2/prices/eth-usd/spot");
  const { data } = (await raw.json()) as Response;
  return data.amount;
}

async function main() {
  const ClientPresence = client.user!.setActivity(`$${await getPrice()}`, {
    type: ActivityTypes.WATCHING,
  });

  console.log(`Activity set to ${ClientPresence.activities[0].name}`);
}

setInterval(() => main(), 30000);

type Response = {
  data: { base: "ETH"; currency: "USD"; amount: string };
};
