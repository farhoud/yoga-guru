import { Card } from "./ui/card"
import { Heading } from "./ui/heading"
import { Text } from "./ui/text"
import { VStack } from "./ui/vstack"
import { HStack } from "./ui/hstack"
import { Button, ButtonIcon, ButtonText } from "./ui/button"
import { Sprout, Leaf, TreePine } from "lucide-react-native"
import { Icon } from "./ui/icon"

const levelIcons = {
  beginner: Sprout,
  intermediate: Leaf,
  advanced: TreePine,
}
export interface CourseListItemProps {
  /**
   * An optional style override useful for padding & margin.
   */
  className?: string
  title: string
  courseType?: string
  schedule: string
  description: string
  level: "beginner" | "intermediate" | "advanced"
  price: number
  capacity: boolean
}

/**
 * Describe your component here
 */
export const CourseListItem = (props: CourseListItemProps) => {
  const { className, title, courseType, schedule, description, level, price, capacity } = props
  return (
    <>
      <Card className="m-2">
        <Heading>{title}</Heading>
        <VStack space="lg" className="py-5">
          <HStack>
            <Text>توضیحات: </Text>
            <Text>{description}</Text>
          </HStack>
          <HStack>
            <Text>زمان: </Text>
            <Text>{schedule}</Text>
          </HStack>
          <HStack>
            <Text>سطح: </Text>
            <Icon as={levelIcons[level]} />
          </HStack>
          <HStack>
            <Text>قیمت: </Text>
            <Text>{price}</Text>
          </HStack>
        </VStack>
        <Button className="bg-primary-500">
          <ButtonText>ثبت نام</ButtonText>
        </Button>
      </Card>
    </>
  )
}
